package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/generated"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/model"
)

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

		err = r.DB.Create(&customer.CustomerAddress).Error
    if err != nil {
        return nil, err
    }

    return &customer, nil
}

func (r *mutationResolver) CreateShippingAddress(ctx context.Context, input model.OrderDropshipperInput) (*model.OrderDropshipper, error) {
	shipping := model.OrderDropshipper {
		Name: input.Name,
		Phone: input.Phone,
		Address1: input.Address1,
		PostalCode: input.PostalCode,
		Country: input.Country,
		Province: input.Province,
		City: input.City,
		SubDistrict: input.SubDistrict,
	}
	err := r.DB.Create(&shipping).Error
	if err != nil {
			return nil, err
	}

	return &shipping, nil
}

func (r *queryResolver) Customers(ctx context.Context) ([]*model.GqlCustomer, error) {
		var customers []*model.GqlCustomer

		err := r.DB.Set("gorm:auto_preload", true).Find(&customers).Error
		if err != nil {
			return nil, err
		}
		return customers, nil
}

func (r *queryResolver) Customer(ctx context.Context, phone string) (*model.GqlCustomer, error) {
		var customer	[]*model.GqlCustomer
		var address		[]*model.GqlCustomerAddresse

		err := r.DB.Set("gorm:auto_preload", true).Where("phone = ?", phone).Find(&customer).Error
		if err != nil {
			return nil, err
		}

		err = r.DB.Set("gorm:auto_preload", true).Where("phone = ?", phone).Find(&address).Error
		if err != nil {
			return nil, err
		}

		data := customer[0]
		data.CustomerAddress = address[0]

		return data, nil
}

func (r *queryResolver) Countries(ctx context.Context) (*model.NewCustomerAddress, error) {
	var address model.NewCustomerAddress
	var countries []*model.Countrie

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name").Find(&countries).Error
	if err != nil {
		return nil, err
	}
	address.Countries = countries
	return &address, nil
}

func (r *queryResolver) Provinces(ctx context.Context, countryID int) (*model.NewCustomerAddress, error) {
	var address model.NewCustomerAddress
	var provinces []*model.Province

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name").Where("country_id = ?", countryID).Find(&provinces).Error
	if err != nil {
		return nil, err
	}

	address.Provinces = provinces
	return &address, nil
}

func (r *queryResolver) Cities(ctx context.Context, provinceID int) (*model.NewCustomerAddress, error) {
	var address model.NewCustomerAddress
	var cities []*model.Citie

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name").Where("province_id = ?", provinceID).Find(&cities).Error
	if err != nil {
		return nil, err
	}

	address.Cities = cities
	return &address, nil
}

func (r *queryResolver) Subdistricts(ctx context.Context, cityID int) (*model.NewCustomerAddress, error) {
	var address model.NewCustomerAddress
	var districts []*model.District
	var subdistricts []*model.SubDistrict

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name").Where("city_id = ?", cityID).Find(&districts).Error
	if err != nil {
		return nil, err
	}

	for _, district := range districts {
		err := r.DB.Set("gorm:auto_preload", true).Select("id, name").Where("district_id = ?", district.ID).Find(&subdistricts).Error
		if err != nil {
			return nil, err
		}
		address.SubDistricts = append(address.SubDistricts,subdistricts...)
	}

	return &address, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }


// =================== resolver support ============================

func mapCustomerAddress(addressInput *model.GqlCustomerAddresseInput) (*model.GqlCustomerAddresse) {
		var address *model.GqlCustomerAddresse
		address = &model.GqlCustomerAddresse {
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
