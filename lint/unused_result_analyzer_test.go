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

package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/tools/analysis"

	"github.com/onflow/cadence-tools/lint"
)

func TestUnusedResultAnalyzer(t *testing.T) {

	t.Parallel()

	t.Run("unary operation", func(t *testing.T) {

		t.Parallel()

		diagnostics := testAnalyzers(t,
			`
			pub contract Test {
				pub fun test() {
					let b = true
					!b
				}
			}
			`,
			lint.UnusedResultAnalyzer,
		)

		require.Equal(
			t,
			[]analysis.Diagnostic{
				{
					Range: ast.Range{
						StartPos: ast.Position{Offset: 68, Line: 5, Column: 5},
						EndPos:   ast.Position{Offset: 69, Line: 5, Column: 6},
					},
					Location: testLocation,
					Category: lint.UpdateCategory,
					Message:  "Unused result of Unary operation.",
				},
			},
			diagnostics,
		)
	})

	t.Run("binary operation", func(t *testing.T) {

		t.Parallel()

		diagnostics := testAnalyzers(t,
			`
			pub contract Test {
				pub fun test() {
					let a = 10
					let b = 20
					a + b
				}
			}
			`,
			lint.UnusedResultAnalyzer,
		)

		require.Equal(
			t,
			[]analysis.Diagnostic{
				{
					Range: ast.Range{
						StartPos: ast.Position{Offset: 82, Line: 6, Column: 5},
						EndPos:   ast.Position{Offset: 86, Line: 6, Column: 9},
					},
					Location: testLocation,
					Category: lint.UpdateCategory,
					Message:  "Unused result of Binary operation.",
				},
			},
			diagnostics,
		)
	})

	// TODO: test function invocation

	t.Run("valid", func(t *testing.T) {

		t.Parallel()

		diagnostics := testAnalyzers(t,
			`
			pub contract Test {
				pub fun test() {
					let a = 10
					let b = 20
					let c = a + b
					var d = true
					d = !d
				}
			}
			`,
			lint.UnusedResultAnalyzer,
		)

		require.Equal(
			t,
			[]analysis.Diagnostic(nil),
			diagnostics,
		)
	})
}
