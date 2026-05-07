//go:build network

package playfab_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/df-mc/go-playfab"
	"github.com/df-mc/go-playfab/catalog"
	"github.com/df-mc/go-playfab/title"
	"github.com/df-mc/go-xsapi/v2"
	"github.com/df-mc/go-xsapi/v2/xal"
	"github.com/df-mc/go-xsapi/v2/xal/sisu"
	"github.com/df-mc/go-xsapi/v2/xal/xasd"
	"github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"
)

func TestSession(t *testing.T) {
	if err := os.MkdirAll(testdataDir, os.ModePerm); err != nil {
		t.Fatalf("error making parent directories for %q: %s", testdataDir, err)
	}
	msa := MinecraftAndroid.TokenSource(context.Background(), msaToken(t, tokenPath))
	t.Cleanup(func() {
		token, err := msa.Token()
		if err != nil {
			t.Errorf("error requesting Microsoft Account token: %s", err)
			return
		}
		b, err := json.Marshal(token)
		if err != nil {
			t.Errorf("error encoding Microsoft Account token for saving: %s", err)
			return
		}
		if err := os.WriteFile(tokenPath, b, os.ModePerm); err != nil {
			t.Errorf("error writing Microsoft Account token to %s: %s", tokenPath, err)
			return
		}
		t.Logf("cleanup: saved Microsoft Account token to %s", tokenPath)
	})

	dt, proofKey := readDevice(t, deviceSnapshotPath)
	deviceSource := xasd.ReuseTokenSource(MinecraftAndroid.Config, dt, proofKey)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		token, err := deviceSource.DeviceToken(ctx)
		if err != nil {
			t.Fatalf("error requesting device token: %s", err)
		}
		writeDevice(t, deviceSnapshotPath, token, deviceSource.ProofKey())
		t.Logf("cleanup: saved device to: %s", deviceSnapshotPath)
	})

	sc := &sisu.SessionConfig{DeviceTokenSource: deviceSource}
	sc.Snapshot = readSnapshot(t, snapshotPath)
	s := MinecraftAndroid.New(msa, sc)
	t.Cleanup(func() {
		cache := s.Snapshot()
		if cache == nil {
			t.Fatal("Session.Snapshot must return non-nil SessionState")
		}
		writeSnapshot(t, snapshotPath, cache)
		t.Logf("cleanup: written session snapshot")
	})

	device, err := s.DeviceToken(tokenContext(t))
	if err != nil {
		t.Fatalf("error requesting XASD token: %s", err)
	}
	t.Logf("device token: %#v", device)
	title, err := s.TitleToken(tokenContext(t))
	if err != nil {
		t.Fatalf("error requesting XAST token: %s", err)
	}
	t.Logf("title token: %#v", title)
	user, err := s.UserToken(tokenContext(t))
	if err != nil {
		t.Fatalf("error requesting XASU token: %s", err)
	}
	t.Logf("user token: %#v", user)

	xsts, err := s.XSTSToken(tokenContext(t), playFabRelyingParty)
	if err != nil {
		t.Fatalf("error requesting XSTS token for %q: %s", playFabRelyingParty, err)
	}
	t.Logf("XSTS token for %q: %#v", playFabRelyingParty, xsts)

	client, err := xsapi.ClientConfig{
		// EnableChat: true,
	}.New(t.Context(), s)
	if err != nil {
		t.Fatalf("error creating API client: %s", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("error closing API client: %s", err)
		}
	})

	pf, err := playfab.LoginWithXbox(t.Context(), MinecraftTitle, client, playfab.ClientConfig{
		CreateAccount: true,
	})
	if err != nil {
		t.Fatalf("login to playfab: %s", err)
	}
	defer pf.Close()

	token, err := pf.MasterPlayerAccount().EntityToken(t.Context())
	if err != nil {
		t.Fatalf("request entity token for master player account: %s", err)
	}
	t.Log(token.Token)

	titlePlayerAccount, err := pf.TitlePlayerAccount().EntityToken(t.Context())
	if err != nil {
		t.Fatalf("request entity token for title player account: %s", err)
	}
	t.Log(titlePlayerAccount.Token)

	t.Log(token.Token == titlePlayerAccount.Token)

	t.Logf("logged in as %s (%s)", client.UserInfo().GamerTag, client.UserInfo().XUID)

	t.Log("Title ID:", title.DisplayClaims.TitleInfo.TitleID)

	searchItems(t, pf, catalog.SearchFilter{
		Term:  "Toy Story",
		Count: 50,
	})
	searchItems(t, pf, catalog.SearchFilter{
		Filter: "(ContentType eq '3PP') OR (ContentType eq '3PP.V2')",
		Count:  50,
	})

	item, err := pf.Catalog().ItemByID(t.Context(), uuid.MustParse("5a56435d-689c-46aa-a685-073750bff4d4"))
	if err != nil {
		t.Fatal(err)
	}
	for _, content := range item.Contents {
		t.Logf("%#v", content)
	}

	t.Logf("%#v", item)

	t.Log(item.Title.Neutral())
	t.Log(item.Title.Lookup(language.MustParse("ja-JP").String()))
}

func searchItems(t testing.TB, pf *playfab.Client, filter catalog.SearchFilter) *catalog.SearchResult {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second*15)
	defer cancel()
	result, err := pf.Catalog().SearchItems(ctx, filter)
	if err != nil {
		t.Fatalf("search items: %s", err)
	}
	t.Logf("result.ContinuationToken = %q", result.ContinuationToken)
	for i, item := range result.Items {
		t.Logf("item #%d: %#v", i, item)
	}
	return result
}

func readSnapshot(t testing.TB, path string) *sisu.Snapshot {
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		t.Fatalf("stat %q: %s", path, err)
	} else if stat.IsDir() {
		t.Fatalf("%q is a directory", path)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading session cache: %s", path)
	}
	var snapshot *sisu.Snapshot
	if err := json.Unmarshal(b, &snapshot); err != nil {
		t.Fatalf("error decoding session cache: %s", err)
	}
	return snapshot
}

func writeSnapshot(t testing.TB, path string, cache *sisu.Snapshot) {
	b, err := json.Marshal(cache)
	if err != nil {
		t.Fatalf("error encoding Snapshot: %s", err)
	}
	if err := os.WriteFile(path, b, os.ModePerm); err != nil {
		t.Fatalf("error writing session snapshot to %s: %s", path, err)
	}
	t.Logf("Session.Snapshot: %s", b)
}

func msaToken(t testing.TB, path string) *oauth2.Token {
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second*15)
		defer cancel()
		da, err := MinecraftAndroid.DeviceAuth(ctx)
		if err != nil {
			t.Fatalf("error requesting device authentication flow: %s", err)
		}
		t.Logf("Sign in to Microsoft Account at %s using the code %s. You have 1 minute to sign in.", da.VerificationURI, da.UserCode)

		ctx, cancel = context.WithTimeout(t.Context(), time.Minute*5)
		defer cancel()
		msa, err := MinecraftAndroid.DeviceAccessToken(ctx, da)
		if err != nil {
			t.Fatalf("error polling device authentication flow for access token: %s", err)
		}
		b, err := json.Marshal(msa)
		if err != nil {
			t.Fatalf("error encoding oauth2 token: %s", err)
		}
		if err := os.WriteFile(path, b, os.ModePerm); err != nil {
			t.Fatalf("error writing oauth2 token to %s: %s", path, err)
		}
		return msa
	} else if err != nil {
		t.Fatalf("stat %q: %s", path, err)
	} else if stat.IsDir() {
		t.Fatalf("%q is a directory", path)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading token cache: %s", err)
	}
	var msa *oauth2.Token
	if err := json.Unmarshal(b, &msa); err != nil {
		t.Fatalf("error decoding oauth2 token from cache: %s", err)
	}
	return msa
}

func tokenContext(t testing.TB) context.Context {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second*15)
	t.Cleanup(cancel)
	return ctx
}

const (
	testdataDir = "testdata"

	playFabRelyingParty = "https://b980a380.minecraft.playfabapi.com/"
)

var (
	snapshotPath = filepath.Join(testdataDir, "session.snapshot")
	tokenPath    = filepath.Join(testdataDir, "msa.token")

	MinecraftTitle = title.Title("20CA2")

	MinecraftAndroid = sisu.Config{
		Config: xal.Config{
			Device: xal.Device{
				Type:    xal.DeviceTypeAndroid,
				Version: "13",
			},
			UserAgent: "XAL Android 2025.04.20250326.000",
			TitleID:   1739947436,
			Sandbox:   "RETAIL",
		},

		ClientID:    "0000000048183522",
		RedirectURI: "ms-xal-0000000048183522://auth",
	}

	serviceConfigID = uuid.MustParse("4fc10100-5f7a-4470-899b-280835760c07")
)

func readDevice(t testing.TB, path string) (*xasd.Token, *ecdsa.PrivateKey) {
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		t.Fatalf("stat %q: %s", path, err)
	} else if stat.IsDir() {
		t.Fatalf("%q is a directory", path)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading device snapshot: %s", err)
	}
	var snapshot *deviceTokenSnapshot
	if err := json.Unmarshal(b, &snapshot); err != nil {
		t.Fatalf("error decoding device snapshot: %s", err)
	}
	return snapshot.DeviceToken, snapshot.ProofKey.Key.(*ecdsa.PrivateKey)
}

func writeDevice(t testing.TB, path string, token *xasd.Token, proofKey *ecdsa.PrivateKey) {
	b, err := json.Marshal(&deviceTokenSnapshot{
		ProofKey: jose.JSONWebKey{
			Key:       proofKey,
			Algorithm: string(jose.ES256),
			Use:       "sig",
		},
		DeviceToken: token,
	})
	if err != nil {
		t.Fatalf("error encoding device token snapshot: %s", err)
	}
	if err := os.WriteFile(path, b, os.ModePerm); err != nil {
		t.Fatalf("error writing device token snapshot: %s", err)
	}
}

type deviceTokenSnapshot struct {
	ProofKey    jose.JSONWebKey
	DeviceToken *xasd.Token
}

var (
	deviceSnapshotPath = filepath.Join(testdataDir, "device.token")
)
