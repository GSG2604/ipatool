package appstore

import (
	"github.com/golang/mock/gomock"
	"github.com/majd/ipatool/pkg/keychain"
	"github.com/majd/ipatool/pkg/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ = Describe("AppStore (Revoke)", Ordered, func() {
	var (
		ctrl         *gomock.Controller
		appstore     AppStore
		mockKeychain *keychain.MockKeychain
	)

	BeforeAll(func() {
		log.Logger = log.Output(log.NewWriter()).Level(zerolog.Disabled)
	})

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockKeychain = keychain.NewMockKeychain(ctrl)
		appstore = NewAppStore(&Args{
			Keychain: mockKeychain,
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("keychain returns error", func() {
		var testErr = errors.New("test error")

		BeforeEach(func() {
			mockKeychain.EXPECT().
				Remove("account").
				Return(testErr)
		})

		It("returns wrapped error", func() {
			err := appstore.Revoke()
			Expect(err).To(MatchError(ContainSubstring(testErr.Error())))
			Expect(err).To(MatchError(ContainSubstring("failed to revoke auth credentials")))
		})
	})

	When("keychain removes item", func() {
		BeforeEach(func() {
			mockKeychain.EXPECT().
				Remove("account").
				Return(nil)
		})

		It("returns data", func() {
			err := appstore.Revoke()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})