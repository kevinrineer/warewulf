package node

import (
	"fmt"
	"github.com/hpcng/warewulf/internal/pkg/errors"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

func New() (nodeYaml, error) {
	var ret nodeYaml

	wwlog.Printf(wwlog.DEBUG, "Opening node configuration file: %s\n", ConfigFile)
	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		fmt.Printf("error reading node configuration file\n")
		return ret, err
	}

	err = yaml.Unmarshal(data, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func (self *nodeYaml) FindAllNodes() ([]NodeInfo, error) {
	var ret []NodeInfo

	for nodename, node := range self.Nodes {
		var n NodeInfo

		n.NetDevs = make(map[string]*NetDevEntry)
		n.SystemOverlay.Set("default")
		n.RuntimeOverlay.Set("default")
		n.Ipxe.Set("default")

		fullname := strings.SplitN(nodename, ".", 2)
		if len(fullname) > 1 {
			n.DomainName.Set(fullname[1])
		}

		if len(node.Profiles) == 0 {
			n.Profiles = []string{"default"}
		} else {
			n.Profiles = node.Profiles
		}

		n.Id.Set(nodename)
		n.Comment.Set(node.Comment)
		n.Vnfs.Set(node.Vnfs)
		n.KernelVersion.Set(node.KernelVersion)
		n.KernelArgs.Set(node.KernelArgs)
		n.DomainName.Set(node.DomainName)
		n.Ipxe.Set(node.Ipxe)
		n.IpmiIpaddr.Set(node.IpmiIpaddr)
		n.IpmiNetmask.Set(node.IpmiNetmask)
		n.IpmiUserName.Set(node.IpmiUserName)
		n.IpmiPassword.Set(node.IpmiPassword)
		n.SystemOverlay.Set(node.SystemOverlay)
		n.RuntimeOverlay.Set(node.RuntimeOverlay)

		for devname, netdev := range node.NetDevs {
			if _, ok := n.NetDevs[devname]; !ok {
				var netdev NetDevEntry
				n.NetDevs[devname] = &netdev
			}

			n.NetDevs[devname].Ipaddr.Set(netdev.Ipaddr)
			n.NetDevs[devname].Netmask.Set(netdev.Netmask)
			n.NetDevs[devname].Hwaddr.Set(netdev.Hwaddr)
			n.NetDevs[devname].Gateway.Set(netdev.Gateway)
			n.NetDevs[devname].Type.Set(netdev.Type)
			n.NetDevs[devname].Default.SetB(netdev.Default)
		}

		for _, p := range n.Profiles {
			if _, ok := self.NodeProfiles[p]; !ok {
				wwlog.Printf(wwlog.WARN, "Profile not found for node '%s': %s\n", nodename, p)
				continue
			}

			wwlog.Printf(wwlog.VERBOSE, "Merging profile into node: %s <- %s\n", nodename, p)

			pstring := fmt.Sprintf("%s", p)

			n.Comment.SetAlt(self.NodeProfiles[p].Comment, pstring)
			n.DomainName.SetAlt(self.NodeProfiles[p].DomainName, pstring)
			n.Vnfs.SetAlt(self.NodeProfiles[p].Vnfs, pstring)
			n.KernelVersion.SetAlt(self.NodeProfiles[p].KernelVersion, pstring)
			n.KernelArgs.SetAlt(self.NodeProfiles[p].KernelArgs, pstring)
			n.Ipxe.SetAlt(self.NodeProfiles[p].Ipxe, pstring)
			n.IpmiIpaddr.SetAlt(self.NodeProfiles[p].IpmiIpaddr, pstring)
			n.IpmiNetmask.SetAlt(self.NodeProfiles[p].IpmiNetmask, pstring)
			n.IpmiUserName.SetAlt(self.NodeProfiles[p].IpmiUserName, pstring)
			n.IpmiPassword.SetAlt(self.NodeProfiles[p].IpmiPassword, pstring)
			n.SystemOverlay.SetAlt(self.NodeProfiles[p].SystemOverlay, pstring)
			n.RuntimeOverlay.SetAlt(self.NodeProfiles[p].RuntimeOverlay, pstring)

			for devname, netdev := range self.NodeProfiles[p].NetDevs {
				if _, ok := n.NetDevs[devname]; !ok {
					var netdev NetDevEntry
					n.NetDevs[devname] = &netdev
				}

				n.NetDevs[devname].Ipaddr.SetAlt(netdev.Ipaddr, pstring)
				n.NetDevs[devname].Netmask.SetAlt(netdev.Netmask, pstring)
				n.NetDevs[devname].Hwaddr.SetAlt(netdev.Hwaddr, pstring)
				n.NetDevs[devname].Gateway.SetAlt(netdev.Gateway, pstring)
				n.NetDevs[devname].Type.SetAlt(netdev.Type, pstring)
				n.NetDevs[devname].Default.SetAltB(netdev.Default, pstring)
			}
		}

		ret = append(ret, n)

	}

	return ret, nil
}

/*
func (self *nodeYaml) FindAllControllers() ([]ControllerConf, error) {
	var ret []ControllerConf

	for controllername, controller := range self.Controllers {
		var c ControllerConf

		c.Id = controllername
		c.Ipaddr = controller.Ipaddr
		c.Comment = controller.Comment
		c.Fqdn = controller.Fqdn

		//TODO: Is there a better way to do this, cause EWWW!
		c.Services = struct {
			Warewulfd struct {
				Port       string
				Secure     bool
				StartCmd   string
				RestartCmd string
				EnableCmd  string
			}
			Dhcp struct {
				Enabled    bool
				Template   string
				RangeStart string
				RangeEnd   string
				ConfigFile string
				StartCmd   string
				RestartCmd string
				EnableCmd  string
			}
			Tftp struct {
				Enabled    bool
				TftpRoot   string
				StartCmd   string
				RestartCmd string
				EnableCmd  string
			}
			Nfs struct {
				Enabled    bool
				Exports    []string
				StartCmd   string
				RestartCmd string
				EnableCmd  string
			}
		}(controller.Services)

		// Validations //

		if c.Ipaddr == "" {
			wwlog.Printf(wwlog.WARN, "Controller IP address is unset: %s\n", c.Id)
		}

		if c.Services.Warewulfd.Port == "" {
			c.Services.Warewulfd.Port = "987"
		}

		// TODO: Validate or die on all inputs

		ret = append(ret, c)
	}
	return ret, nil
}
*/

func (self *nodeYaml) FindAllProfiles() ([]NodeInfo, error) {
	var ret []NodeInfo

	for name, profile := range self.NodeProfiles {
		var p NodeInfo

		p.Id.Set(name)
		p.Comment.Set(profile.Comment)
		p.Vnfs.Set(profile.Vnfs)
		p.Ipxe.Set(profile.Ipxe)
		p.KernelVersion.Set(profile.KernelVersion)
		p.KernelArgs.Set(profile.KernelArgs)
		p.IpmiNetmask.Set(profile.IpmiNetmask)
		p.IpmiUserName.Set(profile.IpmiUserName)
		p.IpmiPassword.Set(profile.IpmiPassword)
		p.RuntimeOverlay.Set(profile.RuntimeOverlay)
		p.SystemOverlay.Set(profile.SystemOverlay)

		for devname, netdev := range profile.NetDevs {
			if _, ok := p.NetDevs[devname]; !ok {
				var netdev NetDevEntry
				p.NetDevs[devname] = &netdev
			}

			p.NetDevs[devname].Ipaddr.Set(netdev.Ipaddr)
			p.NetDevs[devname].Netmask.Set(netdev.Netmask)
			p.NetDevs[devname].Hwaddr.Set(netdev.Hwaddr)
			p.NetDevs[devname].Gateway.Set(netdev.Gateway)
			p.NetDevs[devname].Type.Set(netdev.Type)
			p.NetDevs[devname].Default.SetB(netdev.Default)
		}

		// TODO: Validate or die on all inputs

		ret = append(ret, p)
	}
	return ret, nil
}

func (self *nodeYaml) FindByHwaddr(hwa string) (NodeInfo, error) {
	var ret NodeInfo

	n, _ := self.FindAllNodes()

	for _, node := range n {
		for _, dev := range node.NetDevs {
			if dev.Hwaddr.Get() == hwa {
				return node, nil
			}
		}
	}

	return ret, errors.New("No nodes found with HW Addr: " + hwa)
}

func (self *nodeYaml) FindByIpaddr(ipaddr string) (NodeInfo, error) {
	var ret NodeInfo

	n, _ := self.FindAllNodes()

	for _, node := range n {
		for _, dev := range node.NetDevs {
			if dev.Ipaddr.Get() == ipaddr {
				return node, nil
			}
		}
	}

	return ret, errors.New("No nodes found with IP Addr: " + ipaddr)
}

func (nodes *nodeYaml) SearchByName(search string) ([]NodeInfo, error) {
	var ret []NodeInfo

	n, _ := nodes.FindAllNodes()

	for _, node := range n {
		b, _ := regexp.MatchString(search, node.Id.Get())
		if b == true {
			ret = append(ret, node)
		}
	}

	return ret, nil
}

func (nodes *nodeYaml) SearchByNameList(searchList []string) ([]NodeInfo, error) {
	var ret []NodeInfo

	n, _ := nodes.FindAllNodes()

	for _, search := range searchList {
		for _, node := range n {
			b, _ := regexp.MatchString(search, node.Id.Get())
			if b == true {
				ret = append(ret, node)
			}
		}
	}

	return ret, nil
}
