package cli

import "testing"

func TestFarmModuleList_Validate(t *testing.T) {
	tests := []struct {
		name    string
		list    FarmModuleList
		wantErr bool
	}{
		{
			name: "duplicate export name",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Source:  "http://rowdy-watcher.info",
						Version: "4.5.6",
						Name:    "solution",
					},
					{
						Source:  "https://knobby-courtroom.net",
						Version: "1.2.3",
						Name:    "solution",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate reference",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Source:  "http://bite-sized-scrap.org",
						Version: "9.4.6",
						Name:    "Wooden",
					},
					{
						Source:  "http://bite-sized-scrap.org",
						Version: "9.4.6",
						Name:    "navigate",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid module list",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Source:  "http://defiant-forum.com",
						Version: "0.0.0",
						Name:    "Concrete",
					},
					{
						Source:  "https://heavy-caviar.name",
						Version: "21.5.4",
						Name:    "synthesizing",
					},
					{
						Source:  "http://cultured-subscription.com",
						Version: "8.8.8",
						Name:    "generating",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.list.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("FarmModuleList.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
