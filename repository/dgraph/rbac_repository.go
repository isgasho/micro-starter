package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/micro-community/auth/db"
	"github.com/micro-community/auth/models"
	"github.com/micro/go-micro/v3/logger"
)

type RbacRepository struct {
}

func NewRBACRepository() *RbacRepository {
	return &RbacRepository{}
}

// AddUser is a single request handler called via client.AddUser or the generated client code
func (e *RbacRepository) AddUser(ctx context.Context, user *models.User) error {
	logger.Infof("Received RbacRepository.AddUser request, ID: %d, Name: %s", user.ID, user.Name)
	//首先查询数据库中是否已有该ID

	drsp, err := db.DDB().QueryExist(user.ID)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil || len(r.Count) < 1 {
		return fmt.Errorf("json unmarshal drsp error: %v", err)
	}

	if r.Count[0].Count > 0 {
		return fmt.Errorf("User %d already exists", user.ID)
	}

	// 创建新User
	p := models.User{
		//Uid:    "_:" + target,
		Type:   "User",
		ID:     user.ID,
		Name:   user.Name,
		Age:    user.Age,
		Gender: user.Gender,
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(p)
	if err != nil {
		logger.Fatal(err)
		return fmt.Errorf("json Marshal error: %v", err)
	}

	mu.SetJson = pb
	_, err = db.DDB().Mutate(pb)
	if err != nil {
		return fmt.Errorf("dgraph Mutate error: %v", err)
	}

	return nil
}

// RemoveUser is a single request handler called via client.RemoveUser or the generated client code
func (e *RbacRepository) RemoveUser(ctx context.Context, user *models.User) error {
	logger.Infof("Received RbacRepository.RemoveUser request, ID: %d", user.ID)
	// 首先查询数据库中是否已有该ID

	drsp, err := db.DDB().QueryExist(user.ID)

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return fmt.Errorf("json unmarshal Root error: %v", err)
	}
	err = db.DDB().BatchDelete(r.UID)
	if err != nil {
		return fmt.Errorf("RemoveUser commit error: %v", err)
	}
	return nil
}

// QueryUserRoles is a single request handler called via client.QueryUserRoles or the generated client code
func (e *RbacRepository) QueryUserRoles(ctx context.Context, role *models.Role) ([]*models.Role, error) {
	logger.Infof("Received RbacRepository.QueryUserRoles request, ID: %d", role.ID)

	targetID := fmt.Sprintf("%d", role.ID)
	//	variables := map[string]string{"$id": targetID}
	q := `query Me($id: string){
		roles(func: type(User)) @filter(eq(person.id, $id)) @normalize {
			role {
				id
			  name
			}
		}
	}`

	drsp, err := db.DDB().QueryWithVar(targetID, q)
	if err != nil {
		return nil, fmt.Errorf("query err: %v", err)
	}
	var roles []models.Role
	err = json.Unmarshal(drsp.Json, &roles)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal roles error: %v", err)
	}
	var resRoles []*models.Role
	for _, role := range roles {
		resRoles = append(resRoles, &models.Role{ID: role.ID, Name: role.Name})
	}
	return resRoles, nil
}

// QueryUserResources is a single request handler called via client.QueryUserResources or the generated client code
func (e *RbacRepository) QueryUserResources(ctx context.Context, user *models.User) ([]*models.Resource, error) {
	logger.Infof("Received RbacRepository.QueryUserResources request, ID: %d", user.ID)

	targetID := fmt.Sprintf("%d", user.ID)
	//variables := map[string]string{"$id1": user.ID}

	q := `query Me($id1: string){
		resources(func: type(User)) @filter(eq(person.id, $id1)) @normalize {
			role {
				resource {
					resource.id
					resource.name
				}
			}
		}
	}`

	drsp, err := db.DDB().QueryWithVar(targetID, q)
	if err != nil {
		return nil, fmt.Errorf("query err: %v", err)
	}
	type Root struct {
		Resources []models.Resource `json:"Resource"`
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal Root error: %v", err)
	}

	//过滤重复
	seen := map[int]bool{}
	resRoles := []*models.Resource{}
	for _, res := range r.Resources {
		if !seen[res.ID] {
			seen[res.ID] = true
			resRoles = append(resRoles, &models.Resource{ID: res.ID, Name: res.Name})
		}
	}
	return resRoles, nil
}

// LinkUserRole is a single request handler called via client.LinkUserRole or the generated client code
func (e *RbacRepository) LinkUserRole(ctx context.Context, user *models.User, role *models.Role) error {
	logger.Info("Received RbacRepository.LinkUserRole request: user: %d, role: %d", user.ID, role.ID)

	// 首先查询user id 和 role 对应的 id
	variables := map[string]string{"$id1": user.ID, "$id2": role.ID}
	q := `query Me($id1: string, $rid2id: string){
		user(func: type(User)) @filter(eq(person.id, $id1)) {
			uid
		}
		role(func: type(Role)) @filter(eq(role.id, $id2)) {
			uid
		}
	}`

	drsp, err := db.DDB().QueryWithVar(q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}
	type Root struct {
		UID1 []models.UID `json:"user"`
		UID2 []models.UID `json:"role"`
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return fmt.Errorf("json unmarshal Root error: %v", err)
	}
	if len(r.UID1) == 0 {
		return fmt.Errorf("id1 <%s> not found", req.Id1)
	}
	if len(r.UID2) == 0 {
		return fmt.Errorf("id2 <%s> not found", req.Id2)
	}

	// link
	mu := &api.Mutation{
		CommitNow: true,
	}

	nq := &api.NQuad{
		Subject:   r.UID1[0].UID,
		Predicate: "role",
		ObjectId:  r.UID2[0].UID,
	}
	mu.Set = []*api.NQuad{nq}
	_, err = db.DDB().Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("LinkUserRole Mutate error: %v", err)
	}

//	rsp.Msg = "OK"
	return nil
}

// UnlinkUserRole is a single request handler called via client.UnlinkUserRole or the generated client code
func (e *RbacRepository) UnlinkUserRole(ctx context.Context, role *models.Role) error {
	logger.Info("Received RbacRepository.UnlinkUserRole request: id1: %s, id2: %s", req.Id1, req.Id2)
	// 首先查询id1 和 id2 对应的 uid
	variables := map[string]string{"$id1": req.Id1, "$id2": req.Id2}
	q := `query Me($id1: string, $id2: string){
		find_id1(func: type(User)) @filter(eq(person.id, $id1)) {
			uid
		}
		find_id2(func: type(Role)) @filter(eq(role.id, $id2)) {
			uid
		}
	}`
	drsp, err := db.DDB().QueryWithVars(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}
	type Root struct {
		UID1 []UID `json:"find_id1"`
		UID2 []UID `json:"find_id2"`
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return fmt.Errorf("json unmarshal Root error: %v", err)
	}
	if len(r.UID1) == 0 {
		return fmt.Errorf("id1 <%s> not found", req.Id1)
	}
	if len(r.UID2) == 0 {
		return fmt.Errorf("id2 <%s> not found", req.Id2)
	}

	// unlink
	mu := &api.Mutation{
		CommitNow: true,
	}

	nq := &api.NQuad{
		Subject:   r.UID1[0].UID,
		Predicate: "role",
		ObjectId:  r.UID2[0].UID,
	}
	mu.Del = []*api.NQuad{nq}
	_, err = db.DDB().Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("UnlinkUserRole Mutate error: %v", err)
	}

	rsp.Msg = "OK"
	return nil
}

// AddRole is a single request handler called via client.AddRole or the generated client code
func (e *RbacRepository) AddRole(ctx context.Context, role *models.Role) error {
	logger.Infof("Received RbacRepository.AddRole request, ID: %s, Name: %s", req.Id, req.Name)
	// 首先查询数据库中是否已有该ID
	variables := map[string]string{"$id1": req.Id}
	q := `query Me($id1: string){
		count(func: type(Role)) @filter(eq(role.id, $id1)) {
			count(uid)
		}
	}`
	drsp, err := db.DDB().QueryWithVars(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}

	type Count struct {
		Count int `json:"count"`
	}

	type Root struct {
		Count []Count `json:"count"`
	}
	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil || len(r.Count) < 1 {
		return fmt.Errorf("json unmarshal drsp error: %v", err)
	}

	if r.Count[0].Count > 0 {
		return fmt.Errorf("Role %s already exists", req.Id)
	}

	// 创建新Role
	role := Role{
		Uid:  "_:" + req.Id,
		Type: "Role",
		ID:   req.Id,
		Name: req.Name,
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(role)
	if err != nil {
		logger.Fatal(err)
		return fmt.Errorf("json Marshal error: %v", err)
	}

	mu.SetJson = pb
	result, err := db.DDB().Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("dgraph Mutate error: %v", err)
	}

	rsp.Msg = fmt.Sprintf("role created, id: %s,  uid: %s", req.Id, result.Uids[req.Id])
	return nil
}

// RemoveRole is a single request handler called via client.RemoveRole or the generated client code
func (e *RbacRepository) RemoveRole(ctx context.Context, role *models.Role) error {
	logger.Infof("Received RbacRepository.RemoveRole request, ID: %s", req.Id)
	// 首先查询数据库中是否已有该ID
	variables := map[string]string{"$id1": req.Id}
	q := `query Me($id1: string){
		find(func: type(Role)) @filter(eq(role.id, $id1)) {
			uid
		}
	}`
	drsp, err := db.DDB().QueryWithVars(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}
	logger.Info(string(drsp.Json))

	type Root struct {
		UID []UID `json:"find"`
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return fmt.Errorf("json unmarshal drsp error: %v", err)
	}

	if len(r.UID) == 0 {
		rsp.Msg = fmt.Sprintf("%s not exists", req.Id)
		return nil
	}

	// mutate multiple items, then commit
	txn := db.DDB()
	for _, uid := range r.UID {
		d := map[string]string{"uid": uid.UID}
		logger.Info(d)
		pb, err := json.Marshal(d)
		if err != nil {
			return err
		}
		mu := &api.Mutation{
			DeleteJson: pb,
		}
		drsp, err = txn.Mutate(ctx, mu)
		if err != nil {
			return fmt.Errorf("txn Mutate error: %v", err)
		}
	}
	err = txn.Commit(ctx)
	if err != nil {
		return fmt.Errorf("RemoveRole commit error: %v", err)
	}
	rsp.Msg = "OK"
	return nil
}

// QueryRoleResources is a single request handler called via client.QueryRoleResources or the generated client code
func (e *RbacRepository) QueryRoleResources(ctx context.Context, resource *models.Resource) error {
	logger.Infof("Received RbacRepository.QueryRoleResources request, ID: %s", req.Id)
	variables := map[string]string{"$id1": req.Id}
	q := `query Me($id1: string){
		find(func: type(Role)) @filter(eq(role.id, $id1)) @normalize {
			resource {
				resource.id: resource.id
				resource.name: resource.name
			}
		}
	}`
	drsp, err := db.DDB().QueryWithVars(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}
	type Root struct {
		Resource []Resource `json:"find"`
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return fmt.Errorf("json unmarshal Root error: %v", err)
	}

	for _, res := range r.Resource {
		rsp.Resources = append(rsp.Resources, &models.Resource{Id: res.ID, Name: res.Name})
	}
	return nil
}

// LinkRoleResource is a single request handler called via client.LinkRoleResource or the generated client code
func (e *RbacRepository) LinkRoleResource(ctx context.Context, resource *models.Resource) error {
	logger.Info("Received RbacRepository.LinkRoleResource request: id1: %s, id2: %s", req.Id1, req.Id2)
	// 首先查询id1 和 id2 对应的 uid
	variables := map[string]string{"$id1": req.Id1, "$id2": req.Id2}
	q := `query Me($id1: string, $id2: string){
		find_id1(func: type(Role)) @filter(eq(role.id, $id1)) {
			uid
		}
		find_id2(func: type(Resource)) @filter(eq(resource.id, $id2)) {
			uid
		}
	}`
	drsp, err := db.DDB().QueryWithVars(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}
	type Root struct {
		UID1 []UID `json:"find_id1"`
		UID2 []UID `json:"find_id2"`
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return fmt.Errorf("json unmarshal Root error: %v", err)
	}
	if len(r.UID1) == 0 {
		return fmt.Errorf("id1 <%s> not found", req.Id1)
	}
	if len(r.UID2) == 0 {
		return fmt.Errorf("id2 <%s> not found", req.Id2)
	}

	// link
	mu := &api.Mutation{
		CommitNow: true,
	}

	nq := &api.NQuad{
		Subject:   r.UID1[0].UID,
		Predicate: "resource",
		ObjectId:  r.UID2[0].UID,
	}
	mu.Set = []*api.NQuad{nq}
	_, err = db.DDB().Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("LinkRoleResource Mutate error: %v", err)
	}

	rsp.Msg = "OK"
	return nil
}

// UnlinkRoleResource is a single request handler called via client.UnlinkRoleResource or the generated client code
func (e *RbacRepository) UnlinkRoleResource(ctx context.Context, resource *models.Resource) error {
	logger.Info("Received RbacRepository.UnlinkRoleResource request: id1: %s, id2: %s", req.Id1, req.Id2)
	// 首先查询id1 和 id2 对应的 uid
	variables := map[string]string{"$id1": req.Id1, "$id2": req.Id2}
	q := `query Me($id1: string, $id2: string){
		find_id1(func: type(Role)) @filter(eq(role.id, $id1)) {
			uid
		}
		find_id2(func: type(Resource)) @filter(eq(resource.id, $id2)) {
			uid
		}
	}`
	drsp, err := db.DDB().QueryWithVars(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}
	type Root struct {
		UID1 []UID `json:"find_id1"`
		UID2 []UID `json:"find_id2"`
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return fmt.Errorf("json unmarshal Root error: %v", err)
	}
	if len(r.UID1) == 0 {
		return fmt.Errorf("id1 <%s> not found", req.Id1)
	}
	if len(r.UID2) == 0 {
		return fmt.Errorf("id2 <%s> not found", req.Id2)
	}

	// unlink
	mu := &api.Mutation{
		CommitNow: true,
	}

	nq := &api.NQuad{
		Subject:   r.UID1[0].UID,
		Predicate: "resource",
		ObjectId:  r.UID2[0].UID,
	}
	mu.Del = []*api.NQuad{nq}
	_, err = db.DDB().Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("UnlinkRoleResource Mutate error: %v", err)
	}

	rsp.Msg = "OK"
	return nil
}

// AddResource is a single request handler called via client.AddResource or the generated client code
func (e *RbacRepository) AddResource(ctx context.Context, resource *models.Resource) error {
	logger.Infof("Received RbacRepository.AddResource request, ID: %s, Name: %s", req.Id, req.Name)
	// 首先查询数据库中是否已有该ID
	variables := map[string]string{"$id1": req.Id}
	q := `query Me($id1: string){
		count(func: type(Resource)) @filter(eq(resource.id, $id1)) {
			count(uid)
		}
	}`
	drsp, err := db.DDB().QueryWithVars(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}

	type Count struct {
		Count int `json:"count"`
	}

	type Root struct {
		Count []Count `json:"count"`
	}
	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil || len(r.Count) < 1 {
		return fmt.Errorf("json unmarshal drsp error: %v", err)
	}

	if r.Count[0].Count > 0 {
		return fmt.Errorf("Resource %s already exists", req.Id)
	}

	// 创建新Resource
	res := Resource{
		Uid:  "_:" + req.Id,
		Type: "Resource",
		ID:   req.Id,
		Name: req.Name,
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(res)
	if err != nil {
		logger.Fatal(err)
		return fmt.Errorf("json Marshal error: %v", err)
	}

	mu.SetJson = pb
	result, err := db.DDB().Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("dgraph Mutate error: %v", err)
	}

	rsp.Msg = fmt.Sprintf("resource created, id: %s,  uid: %s", req.Id, result.Uids[req.Id])
	return nil
}

// RemoveResource is a single request handler called via client.RemoveResource or the generated client code
func (e *RbacRepository) RemoveResource(ctx context.Context, resource *models.Resource) error {
	logger.Infof("Received RbacRepository.RemoveResource request, ID: %s", req.Id)
	// 首先查询数据库中是否已有该ID
	variables := map[string]string{"$id1": req.Id}
	q := `query Me($id1: string){
		find(func: type(Resource)) @filter(eq(resource.id, $id1)) {
			uid
		}
	}`
	drsp, err := db.DDB().QueryWithVars(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("query err: %v", err)
	}
	logger.Info(string(drsp.Json))

	type Root struct {
		UID []UID `json:"find"`
	}

	var r Root
	err = json.Unmarshal(drsp.Json, &r)
	if err != nil {
		return fmt.Errorf("json unmarshal drsp error: %v", err)
	}

	if len(r.UID) == 0 {
		rsp.Msg = fmt.Sprintf("%s not exists", req.Id)
		return nil
	}

	// mutate multiple items, then commit
	txn := db.DDB()
	for _, uid := range r.UID {
		d := map[string]string{"uid": uid.UID}
		logger.Info(d)
		pb, err := json.Marshal(d)
		if err != nil {
			return err
		}
		mu := &api.Mutation{
			DeleteJson: pb,
		}
		drsp, err = txn.Mutate(ctx, mu)
		if err != nil {
			return fmt.Errorf("txn Mutate error: %v", err)
		}
	}
	err = txn.Commit(ctx)
	if err != nil {
		return fmt.Errorf("RemoveResource commit error: %v", err)
	}
	rsp.Msg = "OK"
	return nil
}
