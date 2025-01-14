//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rego

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

// The query below should meet the following requirements:
// * Provides no Bindings. Do not use a query that sets a variable, e.g. x := data.signature.allow
// * Queries for a single value.
const QUERY = "data.signature.allow"

func ValidateJSON(jsonBody []byte, entrypoints []string) []error {
	ctx := context.Background()

	r := rego.New(
		rego.Query(QUERY),
		rego.Load(entrypoints, nil))

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return []error{err}
	}

	var input interface{}
	dec := json.NewDecoder(bytes.NewBuffer(jsonBody))
	dec.UseNumber()
	if err := dec.Decode(&input); err != nil {
		return []error{err}
	}

	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return []error{err}
	}

	// Ensure the resultset contains a single result where the Expression contains a single value
	// which is true and there are no Bindings.
	if rs.Allowed() {
		return nil
	}

	var errs []error
	for _, result := range rs {
		for _, expression := range result.Expressions {
			errs = append(errs, fmt.Errorf("expression value, %v, is not true", expression))
		}
	}

	// When rs.Allowed() is not true and len(rs) is 0, the result is undefined. This is a policy
	// check failure.
	if len(errs) == 0 {
		errs = append(errs, fmt.Errorf("result is undefined for query '%s'", QUERY))
	}
	return errs
}
