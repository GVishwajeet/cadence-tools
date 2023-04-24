/*
 * Cadence-lint - The Cadence linter
 *
 * Copyright 2019-2022 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lint

import (
	"fmt"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/tools/analysis"
)

func getDiagnostic(location common.Location, element ast.Element, op string) analysis.Diagnostic {
	return analysis.Diagnostic{
		Location: location,
		Range:    ast.NewRangeFromPositioned(nil, element),
		Category: UpdateCategory,
		Message:  "Unused result of " + op + ".",
	}
}

var UnusedResultAnalyzer = (func() *analysis.Analyzer {

	elementFilter := []ast.Element{
		(*ast.ExpressionStatement)(nil),
	}
	return &analysis.Analyzer{
		Description: "Detects expressions with unused results.",
		Requires: []*analysis.Analyzer{
			analysis.InspectorAnalyzer,
		},
		Run: func(pass *analysis.Pass) interface{} {
			inspector := pass.ResultOf[analysis.InspectorAnalyzer].(*ast.Inspector)

			location := pass.Program.Location
			report := pass.Report

			inspector.Preorder(
				elementFilter,
				func(element ast.Element) {
					expressionStatement, ok := element.(*ast.ExpressionStatement)
					if !ok {
						return
					}

					switch fmt.Sprint(expressionStatement.Expression.ElementType()) {
					case "ElementTypeUnaryExpression":
						report(getDiagnostic(location, element, "Unary operation"))
					case "ElementTypeBinaryExpression":
						report(getDiagnostic(location, element, "Binary operation"))
					case "ElementTypeInvocationExpression":
						// TODO: check function invocation
						fmt.Println(expressionStatement, "ElementTypeInvocationExpression")
						return
					default:
						return
					}
				},
			)
			return nil
		},
	}
})()

func init() {
	RegisterAnalyzer(
		"unused-result",
		UnusedResultAnalyzer,
	)
}
