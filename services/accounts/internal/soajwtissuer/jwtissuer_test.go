package soajwtissuer

import (
	"crypto/ed25519"
	"crypto/rand"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwtIssueAndVerifySimple(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	issuer := New(priv)
	verifier := soajwt.NewEd25519Verifier(pub)

	for i := 0; i < 10; i++ {
		profileId := uuid.New().String()
		jwt, err := issuer.Issue(PersonalData{
			AccountId: i,
			ProfileId: profileId,
		}, time.Hour)

		require.NoError(t, err, "error while issuing token number %d", i)

		token, err := verifier.Verify(jwt)
		require.NoError(t, err, "error while veryfing token number %d", i)

		assert.Equal(t, i, token.AccountId)
		assert.Equal(t, profileId, token.Subject)
	}
}

func TestJwtIssueAndVerifyConcurrent(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	issuer := New(priv)
	verifier := soajwt.NewEd25519Verifier(pub)

	wg := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		num := i
		wg.Add(1)
		go func() {
			profileId := uuid.New().String()
			jwt, err := issuer.Issue(PersonalData{
				AccountId: num,
				ProfileId: profileId,
			}, time.Hour)

			require.NoError(t, err, "error while issuing token number %d", i)

			token, err := verifier.Verify(jwt)
			require.NoError(t, err, "error while veryfing token number %d", i)

			assert.Equal(t, i, token.AccountId)
			assert.Equal(t, profileId, token.Subject)
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestJwtIssueAndVerifyTimeout(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	issuer := New(priv)
	verifier := soajwt.NewEd25519Verifier(pub)

	profileId := uuid.New().String()
	jwt, err := issuer.Issue(PersonalData{
		AccountId: 1,
		ProfileId: profileId,
	}, 5*time.Second)

	require.NoError(t, err, "error while issuing token")

	time.Sleep(6 * time.Second)

	_, err = verifier.Verify(jwt)
	require.Error(t, err, "expired token verified successfully")
}
