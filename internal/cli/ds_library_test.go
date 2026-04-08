package cli

import (
	"strings"
	"testing"
)

func TestDSLibraryImport_SingleKey(t *testing.T) {
	out, err := executeCmd("ds", "library", "import", "abc123def")
	if err != nil {
		t.Fatalf("ds library import: %v", err)
	}
	must := []string{
		`importComponentByKeyAsync("abc123def")`,
		"createInstance()",
		"createdNodeIds",
	}
	for _, s := range must {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q\n--- output ---\n%s", s, out)
		}
	}
}

func TestDSLibraryImport_MultipleKeys(t *testing.T) {
	out, err := executeCmd("ds", "library", "import", "key1", "key2", "key3")
	if err != nil {
		t.Fatalf("ds library import multi: %v", err)
	}
	for _, k := range []string{"key1", "key2", "key3"} {
		if !strings.Contains(out, `importComponentByKeyAsync("`+k+`")`) {
			t.Errorf("output missing import for %s\n%s", k, out)
		}
	}
	if !strings.Contains(out, "inst0.id") || !strings.Contains(out, "inst1.id") || !strings.Contains(out, "inst2.id") {
		t.Errorf("output should return all 3 instance IDs\n%s", out)
	}
}

func TestDSLibraryImport_WithName(t *testing.T) {
	out, err := executeCmd("ds", "library", "import", "abc", "--name", "Hero Component")
	if err != nil {
		t.Fatalf("ds library import --name: %v", err)
	}
	if !strings.Contains(out, `inst0.name = "Hero Component"`) {
		t.Errorf("output missing name assignment\n%s", out)
	}
}

func TestDSLibraryImport_WithParent(t *testing.T) {
	out, err := executeCmd("ds", "library", "import", "abc", "--parent", "123:456")
	if err != nil {
		t.Fatalf("ds library import --parent: %v", err)
	}
	if !strings.Contains(out, `figma.getNodeByIdAsync("123:456")`) {
		t.Errorf("output missing parent lookup\n%s", out)
	}
	if !strings.Contains(out, "appendChild(inst0)") {
		t.Errorf("output missing appendChild\n%s", out)
	}
}

func TestDSLibraryImportSet(t *testing.T) {
	out, err := executeCmd("ds", "library", "import-set", "setKey", "--variant", "Size=Large,State=Default")
	if err != nil {
		t.Fatalf("ds library import-set: %v", err)
	}
	must := []string{
		`importComponentSetByKeyAsync("setKey")`,
		`p["Size"]==="Large"`,
		`p["State"]==="Default"`,
		"createInstance()",
		"createdNodeIds",
	}
	for _, s := range must {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q\n--- output ---\n%s", s, out)
		}
	}
}

func TestDSLibraryImportStyle(t *testing.T) {
	out, err := executeCmd("ds", "library", "import-style", "styleKey")
	if err != nil {
		t.Fatalf("ds library import-style: %v", err)
	}
	if !strings.Contains(out, `importStyleByKeyAsync("styleKey")`) {
		t.Errorf("output missing importStyleByKeyAsync\n%s", out)
	}
	if !strings.Contains(out, "return style.id") {
		t.Errorf("output missing return\n%s", out)
	}
}

func TestDSLibraryImportStyle_WithApply(t *testing.T) {
	out, err := executeCmd("ds", "library", "import-style", "sKey", "--apply", "100:200")
	if err != nil {
		t.Fatalf("ds library import-style --apply: %v", err)
	}
	if !strings.Contains(out, `figma.getNodeByIdAsync("100:200")`) {
		t.Errorf("output missing node lookup\n%s", out)
	}
	if !strings.Contains(out, "fillStyleId") {
		t.Errorf("output missing style application\n%s", out)
	}
}

func TestDSLibraryVariables(t *testing.T) {
	out, err := executeCmd("ds", "library", "variables")
	if err != nil {
		t.Fatalf("ds library variables: %v", err)
	}
	if !strings.Contains(out, "getAvailableLibraryVariableCollectionsAsync") {
		t.Errorf("output missing teamLibrary call\n%s", out)
	}
}

func TestDSLibraryVariables_Collection(t *testing.T) {
	out, err := executeCmd("ds", "library", "variables", "--collection", "col123")
	if err != nil {
		t.Fatalf("ds library variables --collection: %v", err)
	}
	if !strings.Contains(out, `getVariablesInLibraryCollectionAsync("col123")`) {
		t.Errorf("output missing collection-specific call\n%s", out)
	}
}

func TestDSLibraryImport_Composable(t *testing.T) {
	out, err := executeCmd("ds", "library", "import", "k1", "--body-only")
	if err != nil {
		t.Fatalf("ds library import --body-only: %v", err)
	}
	if strings.Contains(out, "createdNodeIds") {
		t.Errorf("body-only should suppress return statement\n%s", out)
	}
	if !strings.Contains(out, "importComponentByKeyAsync") {
		t.Errorf("body-only should still have import call\n%s", out)
	}
}

func TestDSLibraryList_RequiresTeamOrFile(t *testing.T) {
	_, err := executeCmd("ds", "library", "list")
	if err == nil {
		t.Error("expected error when neither --team nor --file is provided")
	}
}
