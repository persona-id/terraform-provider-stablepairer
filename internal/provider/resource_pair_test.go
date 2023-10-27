// Copyright (c) Persona
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccResourcePair(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stablepairer": providerserver.NewProtocol6WithError(New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				resource "stablepairer_pair" "test" {
					keys   = ["a", "b", "c"]
					values = ["1", "2", "3"]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.#", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.0", "a"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.1", "b"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.2", "c"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.#", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.0", "1"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.1", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.2", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.%", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.a", "1"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.b", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.c", "3"),
				),
			},
			{
				Config: `
				resource "stablepairer_pair" "test" {
					keys   = ["a", "b", "c"]
					values = ["1", "2", "3"]
				}
				`,
				PlanOnly: true,
			},
			{
				Config: `
				resource "stablepairer_pair" "test" {
					keys   = ["a", "c"]
					values = ["1", "2", "3"]
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						ExpectResultBeforeAfter{
							Before: map[string]string{
								"a": "1",
								"b": "2",
								"c": "3",
							},
							After: map[string]string{
								"a": "1",
								"c": "3",
							},
						},
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.#", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.0", "a"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.1", "c"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.#", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.0", "1"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.1", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.2", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.%", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.a", "1"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.c", "3"),
				),
			},
			{
				Config: `
				resource "stablepairer_pair" "test" {
					keys   = ["a", "c", "e"]
					values = ["1", "2", "3"]
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						ExpectResultBeforeAfter{
							Before: map[string]string{
								"a": "1",
								"c": "3",
							},
							After: map[string]string{
								"a": "1",
								"c": "3",
								"e": "2",
							},
						},
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.#", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.0", "a"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.1", "c"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.2", "e"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.#", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.0", "1"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.1", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.2", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.%", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.a", "1"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.c", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.e", "2"),
				),
			},
		},
	})
}

func TestAccResourcePairEmptyStart(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stablepairer": providerserver.NewProtocol6WithError(New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				resource "stablepairer_pair" "test" {
					keys   = ["a", "b", "c"]
					values = []
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.#", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.0", "a"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.1", "b"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.2", "c"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.#", "0"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.%", "0"),
				),
			},
			{
				Config: `
				resource "stablepairer_pair" "test" {
					keys   = ["a", "b", "c"]
					values = []
				}
				`,
				PlanOnly: true,
			},
			{
				Config: `
				resource "stablepairer_pair" "test" {
					keys   = ["a", "b", "c"]
					values = ["1", "2"]
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						ExpectResultBeforeAfter{
							Before: map[string]string{},
							After: map[string]string{
								"a": "1",
								"b": "2",
							},
						},
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.#", "3"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.0", "a"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.1", "b"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "keys.2", "c"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.#", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.0", "1"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "values.1", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.%", "2"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.a", "1"),
					resource.TestCheckResourceAttr("stablepairer_pair.test", "result.b", "2"),
				),
			},
		},
	})
}

func TestInternalPairStable(t *testing.T) {
	var tests = []struct {
		keys, values   []basetypes.StringValue
		startingResult map[string]string
		endResult      basetypes.MapValue
	}{
		// All Keys & Values Known
		///////////////////////////////////////////////////////////////////////////
		// empty start
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("2"),
				"c": basetypes.NewStringValue("3"),
			}),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
			},
			startingResult: map[string]string{},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("2"),
			}),
		},
		// stable
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
			}),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
			},
			startingResult: map[string]string{
				"a": "1",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"c": basetypes.NewStringValue("2"),
			}),
		},
		// stable - addition
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringValue("d"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
				"d": basetypes.NewStringValue("4"),
			}),
		},
		// stable - removal
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
			}),
		},
		// stable - addition & removal
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringValue("e"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
				"e": basetypes.NewStringValue("1"),
			}),
		},
		// stable - addition & removal at max
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringValue("e"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
				"d": "4",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
				"e": basetypes.NewStringValue("4"),
			}),
		},
		// stable - addition & removal over max
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringValue("e"),
				basetypes.NewStringValue("f"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
				"d": "4",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
				"e": basetypes.NewStringValue("4"),
			}),
		},
		// Keys and / or Values Known
		///////////////////////////////////////////////////////////////////////////
		// empty start
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringUnknown(),
				basetypes.NewStringValue("c"),
				basetypes.NewStringValue("d"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("2"),
				"c": basetypes.NewStringValue("3"),
				"d": basetypes.NewStringValue("4"),
			}),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringUnknown(),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("2"),
				"c": basetypes.NewStringValue("4"),
			}),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringUnknown(),
			},
			startingResult: map[string]string{},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("2"),
				"c": basetypes.NewStringUnknown(),
			}),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringUnknown(),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringUnknown(),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{},
			endResult:      basetypes.NewMapUnknown(types.StringType),
		},
		// stable
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringUnknown(),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
			}),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringUnknown(),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
			},
			startingResult: map[string]string{
				"a": "1",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"c": basetypes.NewStringValue("2"),
			}),
		},
		// stable - addition
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringUnknown(),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringValue("4"),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapUnknown(types.StringType),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringUnknown(),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
			}),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringValue("d"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringUnknown(),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
				"d": basetypes.NewStringUnknown(),
			}),
		},
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringUnknown(),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"a": basetypes.NewStringValue("1"),
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
			}),
		},
		// stable - addition & removal
		{
			keys: []basetypes.StringValue{
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
				basetypes.NewStringValue("e"),
			},
			values: []basetypes.StringValue{
				basetypes.NewStringValue("1"),
				basetypes.NewStringValue("2"),
				basetypes.NewStringValue("3"),
				basetypes.NewStringUnknown(),
			},
			startingResult: map[string]string{
				"a": "1",
				"b": "3",
				"c": "2",
			},
			endResult: basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{
				"b": basetypes.NewStringValue("3"),
				"c": basetypes.NewStringValue("2"),
				"e": basetypes.NewStringValue("1"),
			}),
		},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%+v,%+v,%+v", test.keys, test.values, test.startingResult)

		t.Run(testname, func(t *testing.T) {
			actualResult := pairStable(test.startingResult, test.keys, test.values)

			if !reflect.DeepEqual(test.endResult, actualResult) {
				t.Errorf("Got %+v, wanted %+v", actualResult, test.endResult)
			}
		})
	}
}

var _ plancheck.PlanCheck = ExpectResultBeforeAfter{}

type ExpectResultBeforeAfter struct {
	Before map[string]string
	After  map[string]string
}

func (pc ExpectResultBeforeAfter) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result error

	for _, rc := range req.Plan.ResourceChanges {
		result = errors.Join(result, convertAndCheck("before", rc.Address, pc.Before, rc.Change.Before))

		result = errors.Join(result, convertAndCheck("after", rc.Address, pc.After, rc.Change.After))
	}

	resp.Error = result
}

func convertAndCheck(name, address string, expectedValue map[string]string, valueInterface interface{}) error {
	var err error

	if value, ok := valueInterface.(map[string]interface{}); ok {
		if result, ok := value["result"]; ok {
			if resultCast, ok := result.(map[string]interface{}); ok {
				resultConverted := make(map[string]string, len(resultCast))

				for key, value := range resultCast {
					if valueCast, ok := value.(string); ok {
						resultConverted[key] = valueCast
					} else {
						err = errors.Join(err, fmt.Errorf("unable to cast %s result for %s at %s", name, address, key))
					}
				}

				if !reflect.DeepEqual(expectedValue, resultConverted) {
					err = errors.Join(err, fmt.Errorf("%s differed for %s, expected %+v but was %+v", name, address, expectedValue, resultConverted))
				}
			} else {
				err = errors.Join(err, fmt.Errorf("unable to cast %s result for %s", name, address))
			}
		} else {
			err = errors.Join(err, fmt.Errorf("unable to read %s result for %s", name, address))
		}
	} else {
		err = errors.Join(err, fmt.Errorf("unable to read %s for %s", name, address))
	}

	return err
}
