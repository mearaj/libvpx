package vp8

type bool_coder_spec struct {
}
type bool_writer struct {
}
type bool_reader struct {
}
type c_bool_coder_spec bool_coder_spec
type c_bool_writer bool_writer
type c_bool_reader bool_reader
type vp8_tree_p *int8
type vp8_token_struct struct {
	Value int
	Len   int
}
type vp8_token vp8_token_struct
