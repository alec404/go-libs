package sliceutil

import (
	"reflect"
	"testing"
)

func TestUniqueSlice(t *testing.T) {
	t.Run("int slice with duplicates", func(t *testing.T) {
		input := []int{1, 2, 2, 3, 1, 4}
		expected := []int{1, 2, 3, 4}
		result := UniqueSlice(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("UniqueSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("string slice with duplicates", func(t *testing.T) {
		input := []string{"a", "b", "a", "c", "b"}
		expected := []string{"a", "b", "c"}
		result := UniqueSlice(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("UniqueSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		expected := []int{}
		result := UniqueSlice(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("UniqueSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("no duplicates", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		expected := []int{1, 2, 3, 4}
		result := UniqueSlice(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("UniqueSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("all same", func(t *testing.T) {
		input := []int{5, 5, 5, 5}
		expected := []int{5}
		result := UniqueSlice(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("UniqueSlice() = %v, want %v", result, expected)
		}
	})
}

func TestRemove(t *testing.T) {
	t.Run("remove by value", func(t *testing.T) {
		input := []interface{}{1, 2, 3, 4}
		predicate := func(v interface{}) bool {
			return v == 2
		}
		result := Remove(input, predicate)
		if len(result) != 3 {
			t.Errorf("Remove() length = %d, want 3", len(result))
		}
		if result[0] != 1 || result[2] != 3 {
			t.Errorf("Remove() = %v, invalid result", result)
		}
	})

	t.Run("remove first match only", func(t *testing.T) {
		input := []interface{}{1, 2, 2, 3}
		predicate := func(v interface{}) bool {
			return v == 2
		}
		result := Remove(input, predicate)
		if len(result) != 3 {
			t.Errorf("Remove() length = %d, want 3", len(result))
		}
	})

	t.Run("no match", func(t *testing.T) {
		input := []interface{}{1, 2, 3}
		predicate := func(v interface{}) bool {
			return v == 5
		}
		result := Remove(input, predicate)
		if !reflect.DeepEqual(result, input) {
			t.Errorf("Remove() = %v, want %v", result, input)
		}
	})
}

func TestRemoveSlice(t *testing.T) {
	t.Run("remove int", func(t *testing.T) {
		input := []int{1, 2, 3, 2, 4}
		expected := []int{1, 3, 2, 4}
		result := RemoveSlice(input, 2)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("remove string", func(t *testing.T) {
		input := []string{"a", "b", "c", "b"}
		expected := []string{"a", "c", "b"}
		result := RemoveSlice(input, "b")
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("target not found", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := []int{1, 2, 3}
		result := RemoveSlice(input, 5)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("remove first element", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := []int{2, 3}
		result := RemoveSlice(input, 1)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("remove last element", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := []int{1, 2}
		result := RemoveSlice(input, 3)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveSlice() = %v, want %v", result, expected)
		}
	})
}

func TestDifferenceSlice(t *testing.T) {
	t.Run("int slices", func(t *testing.T) {
		s1 := []int{1, 2, 3, 4, 5}
		s2 := []int{4, 5, 6, 7, 8}
		expected := []int{6, 7, 8}
		result := DifferenceSlice(s1, s2)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("DifferenceSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("string slices", func(t *testing.T) {
		s1 := []string{"a", "b", "c"}
		s2 := []string{"b", "c", "d", "e"}
		expected := []string{"d", "e"}
		result := DifferenceSlice(s1, s2)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("DifferenceSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("no difference", func(t *testing.T) {
		s1 := []int{1, 2, 3}
		s2 := []int{1, 2, 3}
		result := DifferenceSlice(s1, s2)
		if len(result) != 0 {
			t.Errorf("DifferenceSlice() = %v, want empty slice", result)
		}
	})

	t.Run("all different", func(t *testing.T) {
		s1 := []int{1, 2, 3}
		s2 := []int{4, 5, 6}
		expected := []int{4, 5, 6}
		result := DifferenceSlice(s1, s2)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("DifferenceSlice() = %v, want %v", result, expected)
		}
	})

	t.Run("empty s2", func(t *testing.T) {
		s1 := []int{1, 2, 3}
		s2 := []int{}
		result := DifferenceSlice(s1, s2)
		if len(result) != 0 {
			t.Errorf("DifferenceSlice() = %v, want empty slice", result)
		}
	})

	t.Run("empty s1", func(t *testing.T) {
		s1 := []int{}
		s2 := []int{1, 2, 3}
		expected := []int{1, 2, 3}
		result := DifferenceSlice(s1, s2)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("DifferenceSlice() = %v, want %v", result, expected)
		}
	})
}

func TestInSlice(t *testing.T) {
	t.Run("int found", func(t *testing.T) {
		if !InSlice(2, []int{1, 2, 3}) {
			t.Error("InSlice() = false, want true")
		}
	})

	t.Run("int not found", func(t *testing.T) {
		if InSlice(5, []int{1, 2, 3}) {
			t.Error("InSlice() = true, want false")
		}
	})

	t.Run("string found", func(t *testing.T) {
		if !InSlice("b", []string{"a", "b", "c"}) {
			t.Error("InSlice() = false, want true")
		}
	})

	t.Run("string not found", func(t *testing.T) {
		if InSlice("d", []string{"a", "b", "c"}) {
			t.Error("InSlice() = true, want false")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		if InSlice(1, []int{}) {
			t.Error("InSlice() = true, want false")
		}
	})
}
