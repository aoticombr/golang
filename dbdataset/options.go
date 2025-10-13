package dbdataset

type OptionsDataset func(*DataSet)

func SetName(name string) OptionsDataset {
	return func(ds *DataSet) {
		ds.Name = name
	}
}
