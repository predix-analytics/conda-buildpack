package integration_test

import (
	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pushing an app a second time", func() {
	var app *cutlass.App

	BeforeEach(func() {
		if cutlass.Cached {
			Skip("but running cached tests")
		}

		if isMinicondaTest {
			Skip("Skipping non-miniconda tests")
		}

		app = cutlass.New(Fixtures("no_deps"))
		app.Buildpacks = []string{"python_buildpack"}
	})

	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	Regexp := `\[.*/python\-[\d\.]+\-linux\-x64\-(cflinuxfs.*-)?[\da-f]+\.tgz\]`
	DownloadRegexp := "Download " + Regexp
	CopyRegexp := "Copy " + Regexp

	It("uses the cache for manifest dependencies", func() {
		PushAppAndConfirm(app)
		Expect(app.Stdout.String()).To(MatchRegexp(DownloadRegexp))
		Expect(app.Stdout.String()).ToNot(MatchRegexp(CopyRegexp))

		app.Stdout.Reset()
		PushAppAndConfirm(app)
		Expect(app.Stdout.String()).To(MatchRegexp(CopyRegexp))
		Expect(app.Stdout.String()).ToNot(MatchRegexp(DownloadRegexp))
	})
})
