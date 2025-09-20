package cmd

import (
	"github.com/spf13/cobra"
	"github.com/toastsandwich/kvstore/internal/certgen"
)

var (
	hosts []string
	ip    string
)

var CertGenCmd = &cobra.Command{
	Use:   "certgen -ip ... [-host ...]",
	Short: "used to generate certificates",
	Long:  "certgen is used to generate the certificates for the kvstore",

	RunE: runCertgen,
}

func runCertgen(cmd *cobra.Command, args []string) error {
	return certgen.Generate(ip, hosts)
}

func init() {
	CertGenCmd.Flags().StringVarP(&ip, "ip", "", "", "use this to set ip for certificate")
	CertGenCmd.Flags().StringSliceVar(&hosts, "hosts", nil, "use this to set hosts for certificate")
}
