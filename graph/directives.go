package graph

// 指令
func NewDirectiveRoot() DirectiveRoot {
	return DirectiveRoot{
		AuthCheck: AuthCheck,
	}
}
