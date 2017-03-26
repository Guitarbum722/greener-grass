package greenergrass

import "testing"

const targetTestVersion = 3

func TestTestVersion(t *testing.T) {
	if testVersion != targetTestVersion {
		t.Fatalf("Found testVersion = %v, want %v.", testVersion, targetTestVersion)
	}
}

func TestSeparateName(t *testing.T) {
	if err := LoadTitleData(); err != nil {
		t.Fatalf("Failed to load the prefix/suffix data and got the following error\n%s", err)
	}
	for _, ntc := range nameTestCases {
		got := New(ntc.input)
		got.SeparateName(" ")
		if got.First != ntc.expected.First {
			t.Fatalf("First Name ::: SeparateName(%q) = %q, want %q", ntc.input, got.First, ntc.expected.First)
		}
		if got.Middle != ntc.expected.Middle {
			t.Fatalf("Middle Name ::: SeparateName(%q) = %q, want %q", ntc.input, got.Middle, ntc.expected.Middle)
		}
		if got.Last != ntc.expected.Last {
			t.Fatalf("Last Name ::: SeparateName(%q) = %q, want %q", ntc.input, got.Last, ntc.expected.Last)
		}
		if got.Prefix != ntc.expected.Prefix {
			t.Fatalf("Prefix ::: SeparateName(%q) = %q, want %q", ntc.input, got.Prefix, ntc.expected.Prefix)
		}
		if got.Suffix != ntc.expected.Suffix {
			t.Fatalf("Suffix ::: SeparateName(%q) = %q, want %q", ntc.input, got.Suffix, ntc.expected.Suffix)
		}
	}
}

func BenchmarkSeparateName(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(nameTestCases[2].input).SeparateName(" ")
	}
}

func BenchmarkLoadTitleData(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadTitleData()
	}
}
