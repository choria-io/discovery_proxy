package client

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/choria-io/discovery_proxy/choria"
	"github.com/choria-io/discovery_proxy/client/operations"
	"github.com/choria-io/discovery_proxy/models"
	"github.com/chzyer/readline"
)

type Sets struct {
	Choria      *choria.Choria
	ProxyClient *DiscoveryProxy
}

func NewSets(c *choria.Choria) (*Sets, error) {
	set := Sets{}

	dclient, err := NewDiscoveryProxyClient(c)
	if err != nil {
		return &set, err
	}

	set.Choria = c
	set.ProxyClient = dclient

	return &set, nil
}

func (s *Sets) List(fn func(sets []string) error) error {
	result, err := s.ProxyClient.Operations.GetSets(operations.NewGetSetsParamsWithTimeout(2 * time.Second))

	if err != nil {
		return err
	}

	sets := []string{}
	for _, w := range result.Payload.Sets {
		sets = append(sets, string(w))
	}

	return fn(sets)
}

func (s *Sets) Get(set *string, discover *bool, fn func(set *models.Set) error) error {
	params := operations.NewGetSetSetParamsWithTimeout(2 * time.Second)
	params.Set = *set
	params.Discover = discover

	result, err := s.ProxyClient.Operations.GetSetSet(params)
	if err != nil {
		return err
	}

	return fn(result.Payload)
}

func (s *Sets) Rm(set *string) error {
	params := operations.NewDeleteSetSetParamsWithTimeout(2 * time.Second)
	params.Set = *set

	_, err := s.ProxyClient.Operations.DeleteSetSet(params)

	return err
}

func (s *Sets) Create(set string, query string) error {
	params := operations.NewPostSetParamsWithTimeout(2 * time.Second)
	params.Set = &models.Set{}
	params.Set.Query = &query
	params.Set.Set = models.Word(set)

	_, err := s.ProxyClient.Operations.PostSet(params)

	return err
}

func (s *Sets) Update(set string, query string) error {
	params := operations.NewPutSetSetParamsWithTimeout(2 * time.Second)
	params.Set = set

	params.NewSet = &models.Set{}
	params.NewSet.Query = &query
	params.NewSet.Set = models.Word(set)

	_, err := s.ProxyClient.Operations.PutSetSet(params)

	return err
}

func (s *Sets) HaveSet(set *string) bool {
	err := s.List(func(result []string) error {
		for _, found := range result {
			if *set == found {
				return nil
			}
		}

		return errors.New("not found")
	})

	if err == nil {
		return true
	}

	return false
}

func (s *Sets) PrintSet(set *string, discover bool) error {
	err := s.Get(set, &discover, func(result *models.Set) error {
		fmt.Printf("Details for the '%s' set\n\n", result.Set)
		fmt.Print("Query:\n\n")
		fmt.Printf("    %s\n\n", *result.Query)

		if discover {
			fmt.Print("Matched Nodes:\n\n")
			s.PrintNodes(result.Nodes)
		}

		fmt.Println("")

		return nil
	})

	return err
}

func (s *Sets) PrintNodes(nodes []string) {
	sort.Strings(nodes)

	choria.SliceGroups(nodes, 3, func(group []string) {
		width := readline.GetScreenWidth()/3 - 6

		format := fmt.Sprintf("   %%-%ds", width)

		for _, node := range group {
			fmt.Printf(format, node)
		}
		fmt.Println()
	})
}

func (s *Sets) ResolvePQL(pql string) ([]string, error) {
	params := operations.NewGetDiscoverParamsWithTimeout(10 * time.Second)
	params.Request = &models.DiscoveryRequest{}
	params.Request.Query = pql

	result, err := s.ProxyClient.Operations.GetDiscover(params)
	if err != nil {
		return []string{}, err
	}

	return result.Payload.Nodes, nil
}
