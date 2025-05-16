package main

import (
	"github.com/mynamespace/xyz/sdk/go/xyz"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myRandomResource, err := xyz.NewRandom(ctx, "myRandomResource", &xyz.RandomArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}
		_, err = xyz.NewRandomComponent(ctx, "myRandomComponent", &xyz.RandomComponentArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}
		ctx.Export("output", pulumi.StringMap{
			"value": myRandomResource.Result,
		})
		return nil
	})
}
