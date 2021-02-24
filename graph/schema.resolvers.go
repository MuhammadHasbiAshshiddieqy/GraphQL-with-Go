package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"
	"regexp"
	"strings"
	"strconv"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/generated"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/model"
)

// Key - ContextKey
type Key struct{}

func (r *mutationResolver) CreateCustomer(ctx context.Context, input model.CreateCustomerInput) (*model.CreateCustomerPayload, error) {
	var customerID []*model.Customer
	var customerAddr []*model.CustomerAddress
	// printContextInternals(ctx, false)
	var profileIDInt *int
	profileIDInt = new(int)

	profileID := ctx.Value(Key{}).(string)
	*profileIDInt = fixProfileID(profileID)
	if (*profileIDInt == 0) {
		return nil, nil
	}

	customer := &model.Customer{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		CompanyName: input.CompanyName,
		Email:       input.Email,
		Phone:       input.Phone,
		ProfileID:	 *profileIDInt,
	}

	customerAddress, err := r.mapCustomerAddress(input.BillingAddress, input.ShippingAddress)
	if err != nil {
		return nil, err
	}

	err = r.DB.Create(&customer).Error
	if err != nil {
		return nil, err
	}

	getCustomerID := "phone = ? AND first_name = ? AND last_name = ? AND email = ? AND company_name = ? AND profile_id = ?"
	err = r.DB.Set("gorm:auto_preload", true).Where(getCustomerID, input.Phone, input.FirstName, input.LastName, input.Email, input.CompanyName, profileIDInt).Find(&customerID).Error
	if err != nil {
		return nil, err
	}

	lenCustomerID := len(customerID)-1
	custID, _ := strconv.Atoi(customerID[lenCustomerID].ID)
	customerAddress.CustomerID = custID
	customer.ID = customerID[lenCustomerID].ID

	err = r.DB.Create(&customerAddress).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Set("gorm:auto_preload", true).Where("customer_id = ?", custID).Find(&customerAddr).Error
	if err != nil {
		return nil, err
	}

	customer.CustomerAddress = customerAddr[0]

	createCustomerPayload := &model.CreateCustomerPayload{
		Customer: customer,
	}

	return createCustomerPayload, nil
}

func (r *queryResolver) Customers(ctx context.Context, search string, limit *int) ([]*model.Customer, error) {
	var customers []*model.Customer
	var address []*model.CustomerAddress
	var profileIDInt *int

	profileIDInt = new(int)

	profileID := ctx.Value(Key{}).(string)
	*profileIDInt = fixProfileID(profileID)
	if (*profileIDInt == 0) {
		return nil, nil
	}

	err := r.DB.Set("gorm:auto_preload", true).Where("profile_id = ? AND (first_name like ? OR last_name like ? OR phone like ? OR lower_full_name like ?)", *profileIDInt, search+"%", search+"%", search+"%", "%"+search+"%").Limit(*limit).Find(&customers).Error
	if err != nil {
		return nil, err
	}

	for _, c := range customers {
		err = r.DB.Set("gorm:auto_preload", true).Where("customer_id = ?", c.ID).Find(&address).Error
		if err != nil {
			return nil, err
		}
		if len(address) > 0 {
			c.CustomerAddress = address[0]
		}
	}

	return customers, nil
}

func (r *queryResolver) Customer(ctx context.Context, id int) (*model.Customer, error) {
	var customer []*model.Customer
	var address []*model.CustomerAddress

	err := r.DB.Set("gorm:auto_preload", true).Where("id = ?", id).Find(&customer).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Set("gorm:auto_preload", true).Where("customer_id = ?", id).Find(&address).Error
	if err != nil {
		return nil, err
	}

	data := customer[0]
	if len(address) > 0 {
		data.CustomerAddress = address[0]
	}

	return data, nil
}

func (r *queryResolver) Countries(ctx context.Context) ([]*model.Country, error) {
	var countries []*model.Country

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon").Find(&countries).Error
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *queryResolver) Provinces(ctx context.Context, countryID int) ([]*model.Province, error) {
	var provinces []*model.Province
	var country *model.Country

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon").Where("country_id = ?", countryID).Find(&provinces).Error
	if err != nil {
		return nil, err
	}

	country, err = r.getCountry(countryID)
	if err != nil {
		return nil, err
	}

	for _, p := range provinces {
		p.Country = country
	}

	return provinces, nil
}

func (r *queryResolver) Cities(ctx context.Context, provinceID int) ([]*model.City, error) {
	var cities []*model.City
	var province *model.Province

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon").Where("province_id = ?", provinceID).Find(&cities).Error
	if err != nil {
		return nil, err
	}

	province, err = r.getProvince(provinceID)
	if err != nil {
		return nil, err
	}
	province.Country, err = r.getCountry(*province.CountryID)
	if err != nil {
		return nil, err
	}

	for _, c := range cities {
		c.Province = province
	}

	return cities, nil
}

func (r *queryResolver) Districts(ctx context.Context, cityID int) ([]*model.District, error) {
	var districts []*model.District
	var city *model.City

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon").Where("city_id = ?", cityID).Find(&districts).Error
	if err != nil {
		return nil, err
	}

	city, err = r.getCity(cityID)
	if err != nil {
		return nil, err
	}
	city.Province, err = r.getProvince(*city.ProvinceID)
	if err != nil {
		return nil, err
	}
	city.Province.Country, err = r.getCountry(*city.Province.CountryID)
	if err != nil {
		return nil, err
	}

	for _, d := range districts {
		d.City = city
	}

	return districts, nil
}

func (r *queryResolver) SubDistricts(ctx context.Context, districtID int) ([]*model.SubDistrict, error) {
	var subdistricts []*model.SubDistrict
	var district *model.District

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon, postal_code, district_id").Where("district_id = ?", districtID).Find(&subdistricts).Error
	if err != nil {
		return nil, err
	}

	district, err = r.getDistrict(districtID)
	if err != nil {
		return nil, err
	}
	district.City, err = r.getCity(*district.CityID)
	if err != nil {
		return nil, err
	}
	district.City.Province, err = r.getProvince(*district.City.ProvinceID)
	if err != nil {
		return nil, err
	}
	district.City.Province.Country, err = r.getCountry(*district.City.Province.CountryID)
	if err != nil {
		return nil, err
	}

	for _, s := range subdistricts {
		s.District = district
	}

	return subdistricts, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.


// ============================================== MUTATION HELPER ========================================================

func (r *mutationResolver) getCountry(countryID int) (*model.Country, error) {
	var country []*model.Country
	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon").Where("id = ?", countryID).Find(&country).Error
	if err != nil {
		return nil, err
	}
	return country[0], nil
}

func (r *mutationResolver) getProvince(provinceID int) (*model.Province, error) {
	var province []*model.Province
	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon, country_id").Where("id = ?", provinceID).Find(&province).Error
	if err != nil {
		return nil, err
	}

	return province[0], nil
}

func (r *mutationResolver) getCity(cityID int) (*model.City, error) {
	var city []*model.City

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon, province_id").Where("id = ?", cityID).Find(&city).Error
	if err != nil {
		return nil, err
	}

	return city[0], nil
}

func (r *mutationResolver) getDistrict(districtID int) (*model.District, error) {
	var district []*model.District

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon, city_id").Where("id = ?", districtID).Find(&district).Error
	if err != nil {
		return nil, err
	}

	return district[0], nil
}

func (r *mutationResolver) getSubDistrictDetail(SubDistrictID int) (*model.SubDistrict, error) {
	var subdistrict []*model.SubDistrict
	var district *model.District

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon, postal_code, district_id").Where("id = ?", SubDistrictID).Find(&subdistrict).Error
	if err != nil {
		return nil, err
	}

	district, err = r.getDistrict(*subdistrict[0].DistrictID)
	if err != nil {
		return nil, err
	}
	district.City, err = r.getCity(*district.CityID)
	if err != nil {
		return nil, err
	}
	district.City.Province, err = r.getProvince(*district.City.ProvinceID)
	if err != nil {
		return nil, err
	}
	district.City.Province.Country, err = r.getCountry(*district.City.Province.CountryID)
	if err != nil {
		return nil, err
	}

	subdistrict[0].District = district

	return subdistrict[0], nil
}

func (r *mutationResolver) mapCustomerAddress(billingInput *model.BillingAddressInput, shippingInput *model.ShippingAddressInput) (*model.CustomerAddress, error) {
	var address *model.CustomerAddress
	var subDistrict *model.SubDistrict
	var billingSubDistrict *model.SubDistrict

	subDistrict, err := r.getSubDistrictDetail(shippingInput.SubDistrictID)
	if err != nil {
		return nil, err
	}

	billingSubDistrict, err = r.getSubDistrictDetail(billingInput.BillingSubDistrictID)
	if err != nil {
		return nil, err
	}

	address = &model.CustomerAddress{
		Name:               shippingInput.Name,
		Phone:              shippingInput.Phone,
		Address:            shippingInput.Address,
		PostalCode:         shippingInput.PostalCode,
		Country:            *subDistrict.District.City.Province.Country.Name,
		Province:           *subDistrict.District.City.Province.Name,
		City:               *subDistrict.District.City.Name,
		District:           subDistrict.District.Name,
		SubDistrict:        *subDistrict.Name,
		BillingName:        billingInput.BillingName,
		BillingPhone:       billingInput.BillingPhone,
		BillingAddress:     billingInput.BillingAddress,
		BillingPostalCode:  billingInput.BillingPostalCode,
		BillingCountry:     *billingSubDistrict.District.City.Province.Country.Name,
		BillingProvince:    *billingSubDistrict.District.City.Province.Name,
		BillingCity:        *billingSubDistrict.District.City.Name,
		BillingDistrict:    billingSubDistrict.District.Name,
		BillingSubDistrict: *billingSubDistrict.Name,
	}

	return address, nil
}

// ============================================== QUERY HELPER ========================================================

func (r *queryResolver) getCountry(countryID int) (*model.Country, error) {
	var country []*model.Country

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon").Where("id = ?", countryID).Find(&country).Error
	if err != nil {
		return nil, err
	}

	return country[0], nil
}
func (r *queryResolver) getProvince(provinceID int) (*model.Province, error) {
	var province []*model.Province

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon, country_id").Where("id = ?", provinceID).Find(&province).Error
	if err != nil {
		return nil, err
	}

	return province[0], nil
}
func (r *queryResolver) getCity(cityID int) (*model.City, error) {
	var city []*model.City

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon, province_id").Where("id = ?", cityID).Find(&city).Error
	if err != nil {
		return nil, err
	}

	return city[0], nil
}
func (r *queryResolver) getDistrict(districtID int) (*model.District, error) {
	var district []*model.District

	err := r.DB.Set("gorm:auto_preload", true).Select("id, name, lat, lon, city_id").Where("id = ?", districtID).Find(&district).Error
	if err != nil {
		return nil, err
	}

	return district[0], nil
}


//========================================= GENERAL HELPER =================================================

func printContextInternals(ctx interface{}, inner bool) {
    contextValues := reflect.ValueOf(ctx).Elem()
    contextKeys := reflect.TypeOf(ctx).Elem()

    if !inner {
        fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
    }

    if contextKeys.Kind() == reflect.Struct {
        for i := 0; i < contextValues.NumField(); i++ {
            reflectValue := contextValues.Field(i)
            reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

            reflectField := contextKeys.Field(i)

            if reflectField.Name == "Context" {
                printContextInternals(reflectValue.Interface(), true)
            } else {
                fmt.Printf("field name: %+v\n", reflectField.Name)
                fmt.Printf("value: %+v\n", reflectValue.Interface())
            }
        }
    } else {
        fmt.Printf("context is empty (int)\n")
    }
}

func fixProfileID(profID string) (int) {
	findProfID := regexp.MustCompile(`"profile_id":?`)
	spltr := regexp.MustCompile(`[,}]`)
	splittedHeader := spltr.Split(profID, -1)

	for _, s := range(splittedHeader) {
		profIdx := findProfID.FindStringIndex(s)
		if (profIdx!=nil) {
			profileID := strings.TrimSpace(s[profIdx[1]:len(s)])
			profileIDInt, _ := strconv.Atoi(profileID)
			return profileIDInt
		}
	}

	return 0
}
