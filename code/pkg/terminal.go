package pkg

type Terminal struct {
}

func (t *Terminal) New() (PowerShell, error) {
	return NewLocalPowerShell("powershell.exe", "-NoExit", "-Command", "-")
}
