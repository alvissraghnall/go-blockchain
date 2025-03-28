package cli

import (
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	// "text/tabwriter"
  "github.com/jedib0t/go-pretty/v6/table"
)

func printTable(data interface{}, fields []string) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetStyle(table.StyleLight)

    headers := make([]interface{}, len(fields))
    for i, field := range fields {
        headers[i] = field
    }
    t.AppendHeader(headers)

    sliceValue := reflect.ValueOf(data)
    if sliceValue.Kind() != reflect.Slice {
        fmt.Println("Error: data must be a slice")
        return
    }

    for i := 0; i < sliceValue.Len(); i++ {
        item := sliceValue.Index(i)
        if item.Kind() == reflect.Ptr {
            item = item.Elem()
        }

        row := make([]interface{}, len(fields))
        for j, field := range fields {
            var fieldValue reflect.Value

            if item.Type() == reflect.TypeOf(CliWallet{}) {
                if field == "DefaultWallet" {
                    fieldValue = item.FieldByName(field)
                } else {
                    walletField := item.FieldByName("Wallet")
                    if walletField.IsNil() {
                        row[j] = "N/A"
                        continue
                    }
                    fieldValue = walletField.Elem().FieldByName(field)
                }
            } else {
                fieldValue = item.FieldByName(field)
            }

            if !fieldValue.IsValid() {
                row[j] = "N/A"
                continue
            }

            var strValue string
            switch fieldValue.Kind() {
            case reflect.Bool:
                if fieldValue.Bool() {
                    strValue = "*"
                } else {
                    strValue = " "
                }
            case reflect.Ptr:
                if fieldValue.IsNil() {
                    strValue = "nil"
                } else {
                    switch field {
                    case "PrivateKey", "PublicKey":
                        strValue = "[redacted]"
                    default:
                        strValue = fmt.Sprintf("%v", fieldValue.Interface())
                    }
                }
            case reflect.Slice:
                if field == "Address" && fieldValue.Len() > 0 {
                    strValue = hex.EncodeToString(fieldValue.Bytes()[:4]) + "..." +
                        hex.EncodeToString(fieldValue.Bytes()[len(fieldValue.Bytes())-4:])
                } else {
                    strValue = fmt.Sprintf("%v", fieldValue.Interface())
                }
            default:
                strValue = fmt.Sprintf("%v", fieldValue.Interface())
            }

            switch field {
            case "Mnemonic":
                if len(strValue) > 0 {
                    strValue = "[redacted]"
                }
            case "PrivateKey", "PublicKey":
                strValue = "[redacted]"
            }

            row[j] = strValue
        }
        t.AppendRow(row)
    }

    t.Render()
}

/**
	printTable([]*Wallet{wallets[0].Wallet, wallets[1].Wallet}, []string{"Alias", "Address"})
*/
