package jsvm

import (
	"fmt"

	"github.com/dop251/goja"
)

type S3MetaData struct {
}

type S3Bucket interface {
	UploadObject(string, interface{}, interface{}) error
	DownloadObject(string, interface{}) (interface{}, error)
	DeleteObject(string, interface{}) error
}

type LocalVMS3Bucket struct {
}

func (o *LocalVMS3Bucket) UploadObject(name string, data interface{}, mData interface{}) error {
	fmt.Println("S3_BUCKET:Upload", name, data, mData)
	return nil
}

func (o *LocalVMS3Bucket) DownloadObject(name string, mData interface{}) (interface{}, error) {
	return nil, nil
}

func (o *LocalVMS3Bucket) DeleteObject(name string, mData interface{}) error {
	return nil
}

func sthreeb_uploadObject(bucket S3Bucket, call goja.FunctionCall) goja.Value {
	// Die Parameter werden abgerufen
	name, data, metaData := call.Arguments[0].String(), call.Arguments[1].Export(), call.Arguments[2].Export()

	// Das Objekt wird hochgeladen/geschrieben
	if err := bucket.UploadObject(name, data, metaData); err != nil {
		panic(err)
	}

	// Der Vorgang wurde ohne Fehler durchgeführt
	return goja.Undefined()
}

func sthreeb_downloadObject(bucket S3Bucket, runtime *goja.Runtime, call goja.FunctionCall) goja.Value {
	// Die Parameter werden abgerufen
	name, metaData := call.Arguments[0].String(), call.Arguments[1].Export()

	// Das Objekt wird hochgeladen/geschrieben
	downloadedObject, err := bucket.DownloadObject(name, metaData)
	if err != nil {
		panic(err)
	}

	// Der Vorgang wurde ohne Fehler durchgeführt
	return runtime.ToValue(downloadedObject)
}

func sthreeb_deleteObject(bucket S3Bucket, call goja.FunctionCall) goja.Value {
	// Die Parameter werden abgerufen
	name, metaData := call.Arguments[0].String(), call.Arguments[1].Export()

	// Das Objekt wird hochgeladen/geschrieben
	if err := bucket.DeleteObject(name, metaData); err != nil {
		panic(err)
	}

	// Der Vorgang wurde ohne Fehler durchgeführt
	return goja.Undefined()
}

func sthreeb_init(bucketNameOrUrl string, runtime *goja.Runtime, call goja.FunctionCall, vm *JsVM) goja.Value {
	_ = call

	// Es wird geprüft ob das Bucket oder die BucketURL zulässig ist
	if !vm.validateS3BucketEndPoint(bucketNameOrUrl) {
		panic(runtime.NewTypeError("Zweites Argument ist keine Funktion"))
	}

	// Die S3 Funktionen werden bereitgestellt
	bucket, err := vm.initS3Bucket(bucketNameOrUrl)
	if err != nil {
		panic(runtime.NewTypeError("Zweites Argument ist keine Funktion"))
	}

	// Das Goja-JS Objekt wird erstellt
	newGoJaJSObject := runtime.NewObject()

	// Die Upload Funktion wird bereitgestellt
	newGoJaJSObject.Set("uploadObject", func(parm goja.FunctionCall) goja.Value { return sthreeb_uploadObject(bucket, parm) })

	// Die Download Funktion wird bereitgestellt
	newGoJaJSObject.Set("downloadObject", func(parm goja.FunctionCall) goja.Value { return sthreeb_downloadObject(bucket, runtime, parm) })

	// Die Delete Funktion wird bereitgesellt
	newGoJaJSObject.Set("deleteObject", func(parm goja.FunctionCall) goja.Value { return sthreeb_deleteObject(bucket, parm) })

	// Das Objekt wird zurückgegeben
	return newGoJaJSObject
}

func sthreeb_base(runtime *goja.Runtime, call goja.FunctionCall, vm *JsVM) goja.Value {
	_ = call
	return runtime.ToValue(func(parms goja.FunctionCall) goja.Value {
		switch parms.Arguments[0].String() {
		case "init":
			return sthreeb_init(parms.Arguments[1].String(), runtime, call, vm)
		default:
			return goja.Undefined()
		}
	})
}
