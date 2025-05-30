package apiprofile

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"gopkg.in/yaml.v3"
)

/*
Returns the formatted list of profiles as string
*/
func ProfileList(ShowOpt *wwapiv1.GetProfileList) (profileList wwapiv1.ProfileList, err error) {
	profileList.Output = []string{}
	nodeDB, err := node.New()
	if err != nil {
		return
	}
	profiles, err := nodeDB.FindAllProfiles()
	if err != nil {
		return
	}
	profiles = node.FilterProfileListByName(profiles, ShowOpt.Profiles)
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Id() < profiles[j].Id()
	})
	if ShowOpt.ShowAll {
		profileList.Output = append(profileList.Output,
			fmt.Sprintf("%s:=:%s:=:%s", "PROFILE", "FIELD", "VALUE"))
		for _, p := range profiles {
			fields := node.GetFieldList(p)
			for _, f := range fields {
				profileList.Output = append(profileList.Output,
					fmt.Sprintf("%s:=:%s:=:%s", p.Id(), f.Field, f.Value))
			}
		}
	} else if ShowOpt.ShowYaml {
		profileMap := make(map[string]node.Profile)
		for _, profile := range profiles {
			profileMap[profile.Id()] = profile
		}

		buf, _ := yaml.Marshal(profileMap)
		profileList.Output = append(profileList.Output, string(buf))
	} else if ShowOpt.ShowJson {
		profileMap := make(map[string]node.Profile)
		for _, profile := range profiles {
			profileMap[profile.Id()] = profile
		}

		buf, _ := json.MarshalIndent(profileMap, "", "  ")
		profileList.Output = append(profileList.Output, string(buf))
	} else {
		profileList.Output = append(profileList.Output,
			fmt.Sprintf("%s:=:%s", "PROFILE NAME", "COMMENT/DESCRIPTION"))

		for _, profile := range profiles {
			profileList.Output = append(profileList.Output,
				fmt.Sprintf("%s:=:%s", profile.Id(), profile.Comment))
		}
	}
	return
}
