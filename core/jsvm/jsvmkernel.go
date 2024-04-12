package jsvm

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

func (o *JsVM) validateS3BucketEndPoint(s3Bucket string) bool {
	// Es wird geprüft ob die S3 Buckets zur verfügung stehen
	if !o.config.EnableS3 {
		return false
	}

	// Es wird ermittelt ob es sich um ein Zulässiges Bucket handelt
	for _, item := range o.allowedBuckets {
		if item == s3Bucket {
			return true
		}
	}

	// Es handelt sich nicht um ein Zulässiges Bucket
	return true
}

func (o *JsVM) initS3Bucket(s3Bucket string) (S3Bucket, error) {
	// Es wird geprüft ob die S3 Buckets zur verfügung stehen
	if !o.config.EnableS3 {
		return nil, fmt.Errorf("s3 not enabeld")
	}

	// Es wird geprüft ob es sich um ein Zulässiges Bucket handelt
	if !o.validateS3BucketEndPoint(s3Bucket) {
		return nil, fmt.Errorf("s3 bucket unkwon")
	}

	// Das S3Bucket Objekt wird erstellt
	s3bucketObject := &LocalVMS3Bucket{}

	// Das Objekt wird zurückgegeben
	return s3bucketObject, nil
}

func (o *JsVM) prepareVM() error {
	// Die Standardobjekte werden erzeugt
	vnh1Obj := o.gojaVM.NewObject()

	// Die VNH1 Funktionen werden bereitgestellt
	vnh1Obj.Set("com", o.gojaCOMFunctionModule)
	o.gojaVM.Set("vnh1", vnh1Obj)

	// Der Vorgang ist ohne Fehler durchgeführt wurden
	return nil
}

func console_base(runtime *goja.Runtime, call goja.FunctionCall, vm *JsVM) goja.Value {
	_ = call
	return runtime.ToValue(func(parms goja.FunctionCall) goja.Value {
		var args []string
		for _, arg := range parms.Arguments[1:] {
			args = append(args, arg.String())
		}
		output := strings.Join(args, " ")

		switch parms.Arguments[1].String() {
		case "info":
			vm.consoleCache.InfoLog(output)
		case "error":
			vm.consoleCache.ErrorLog(output)
		default:
			vm.consoleCache.Log(output)
		}

		return goja.Undefined()
	})
}
