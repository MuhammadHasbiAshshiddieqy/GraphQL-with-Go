package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/generated"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/model"
)

func mapCustomerAddress(addressInput *model.GqlCustomerAddresseInput) (*model.GqlCustomerAddresse) {
		address := &model.GqlCustomerAddresse {
				Name: addressInput.Name,
				Phone: addressInput.Phone,
				Address: addressInput.Address,
				PostalCode: addressInput.PostalCode,
				Country: addressInput.Country,
				Province: addressInput.Province,
				City: addressInput.City,
				SubDistrict: addressInput.SubDistrict,
				BillingName: addressInput.BillingName,
				BillingPhone: addressInput.BillingPhone,
				BillingAddress: addressInput.BillingAddress,
				BillingPostalCode: addressInput.BillingPostalCode,
				BillingCountry: addressInput.BillingCountry,
				BillingProvince: addressInput.BillingProvince,
				BillingCity: addressInput.BillingCity,
				BillingSubDistrict: addressInput.BillingSubDistrict,
		}
    return address
}

func (r *mutationResolver) CreateCustomer(ctx context.Context, input model.GqlCustomerInput) (*model.GqlCustomer, error) {
		customer := model.GqlCustomer {
				FirstName: input.FirstName,
				LastName: input.LastName,
				CompanyName: input.CompanyName,
				Email: input.Email,
				Phone: input.Phone,
				CustomerAddress: mapCustomerAddress(input.CustomerAddress),
    }
		err := r.DB.Create(&customer).Error
    if err != nil {
        return nil, err
    }
    return &customer, nil
}

func (r *queryResolver) Customers(ctx context.Context) ([]*model.GqlCustomer, error) {
		var customers []*model.GqlCustomer

		err := r.DB.Set("gorm:auto_preload", true).Find(&customers).Error
		if err != nil {
			return nil, err
		}
		return customers, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
