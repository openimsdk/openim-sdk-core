package process

func AddConditions(p *Process, condition ...bool) {
	p.RunConditions = append(p.RunConditions, condition...)
}
