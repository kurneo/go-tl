package filesystem

type LocalDisk struct {
	LocalDiskContract
}

func NewLocalDisk(driver DriverContract) LocalDiskContract {
	return LocalDisk{
		disk{
			driver: driver,
		},
	}
}
