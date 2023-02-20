package main

import (
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"os"
	"strings"

	dcclient "fybrik.io/fybrik/pkg/connectors/datacatalog/clients"
	"fybrik.io/fybrik/pkg/model/datacatalog"
	"fybrik.io/fybrik/pkg/taxonomy/validate"
)

var version string

const (
	requestJsonOption         = "request"
	requestOperationOption    = "operation"
	credentialPathOption      = "creds"
	catalogconnectorUrlOption = "url"
	datasetIDOption           = "dataset"
)

var (
	requestFile         string
	requestOperation    string
	credentialPath      string
	catalogconnectorUrl string
	datasetID           string
)

var DataCatalogTaxonomy = "resources/taxonomy/datacatalog.json#/definitions/GetAssetResponse"

func newDataCatalog() (dcclient.DataCatalog, error) {
	providerName := "egeria"
	return dcclient.NewDataCatalog(
		providerName,
		catalogconnectorUrl)
}

func ValidateAssetResponse(response *datacatalog.GetAssetResponse) error {
	var allErrs []*field.Error
	taxonomyFile := DataCatalogTaxonomy

	// Convert GetAssetRequest Go struct to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println("responseJSON:" + string(responseJSON))

	// Validate Fybrik module against taxonomy
	allErrs, err = validate.TaxonomyCheck(responseJSON, taxonomyFile)
	if err != nil {
		return err
	}

	// Return any error
	if len(allErrs) == 0 {
		return nil
	}

	return errors.New("all Err is not null")
}

func handleRead() error {
	// Initialize DataCatalog interface
	catalog, err := newDataCatalog()
	if err != nil {
		return errors.Wrap(err, "unable to create data catalog facade")
	}
	defer catalog.Close()

	// Open our jsonFile
	jsonFile, err := os.Open(requestFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		return errors.Wrap(err, "error opening "+requestFile)
	}
	fmt.Println("Successfully Opened " + requestFile)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dataCatalogReq datacatalog.GetAssetRequest
	json.Unmarshal(byteValue, &dataCatalogReq)
	var response *datacatalog.GetAssetResponse
	fmt.Println(requestFile)

	if response, err = catalog.GetAssetInfo(&dataCatalogReq, credentialPath); err != nil {
		return errors.Wrap(err, "failed to receive the catalog connector response")
	}
	err = ValidateAssetResponse(response)
	if err != nil {
		return errors.Wrap(err, "failed to validate the catalog connector response")
	}
	fmt.Println("RESPONSE VALIDATION PASS")
	return nil

}

// RootCmd defines the root cli command
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "catalog-fake-client",
		Short:         "Data catalog fake client",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       strings.TrimSpace(version),
		RunE: func(cmd *cobra.Command, args []string) error {
			if requestOperation == "read" {
				return handleRead()
			}
			return errors.New("Unsupported operation")
		},
	}
	cmd.PersistentFlags().StringVar(&requestFile, requestJsonOption, "resources/read-request.json", "Json file containing the data catalog request")
	cmd.PersistentFlags().StringVar(&requestOperation, requestOperationOption, "read", "Request operation")
	cmd.PersistentFlags().StringVar(&credentialPath, credentialPathOption, "cccc", "Credential path")
	cmd.PersistentFlags().StringVar(&catalogconnectorUrl, catalogconnectorUrlOption, "https://localhost:8888", "Catalog connector Url")
	cmd.PersistentFlags().StringVar(&datasetID, datasetIDOption, "qqq", "Dataset ID")
	//cmd.MarkFlagsRequiredTogether(requestJsonOption, requestOperationOption, credentialPathOption, catalogconnectorUrlOption, datasetIDOption)

	return cmd
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
