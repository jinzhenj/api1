package api1

func (s *Schema) MergeGroupIfaces() {
	for i, g := range s.Groups {
		var ifaces []Iface
		for _, new := range g.Ifaces {
			var merged bool
			for j, old := range ifaces {
				if new.Name == old.Name {
					ifaces[j] = mergeIfaces(old, new)
					merged = true
					break
				}
			}
			if !merged {
				ifaces = append(ifaces, new)
			}
		}
		g.Ifaces = ifaces
		s.Groups[i] = g
	}
}

func mergeIfaces(a Iface, b Iface) Iface {
	a.Comments = append(a.Comments, b.Comments...)
	a.PostComments = append(a.PostComments, b.PostComments...)
	for key, val := range b.SemComments {
		a.HasComments.AddSemComment(key, val)
	}
	a.Funs = append(a.Funs, b.Funs...)
	return a
}
