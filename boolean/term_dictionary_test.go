package boolean

import (
	"reflect"
	"testing"
)


func TestTermDictionary_intersectMultiple(t *testing.T) {
	type fields struct {
		Terms []Term
	}
	type args struct {
		list [][]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			name:   "No result",
			fields: fields{},
			args:   args{ [][]int{{2,4,5},{},{3,9,1,0},{3},{2,3}} },
			want:   []int{},
		},
		{
			name:   "One list",
			fields: fields{},
			args:   args{ [][]int{ {2,4,5} } },
			want:   []int{2,4,5},
		},
		{
			name:   "Empty list",
			fields: fields{},
			args:   args{ [][]int{} },
			want:   nil,
		},
		{
			name:   "Merged list",
			fields: fields{},
			args:   args{ [][]int{ {2,4,5}, {4,5,6} } },
			want:   []int{4,5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &TermDictionary{
				Terms: tt.fields.Terms,
			}
			if got := this.intersectMultiple(tt.args.list...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TermDictionary.intersectMultiple() = %v, want %v", got, tt.want)
			}
		})
	}
}
