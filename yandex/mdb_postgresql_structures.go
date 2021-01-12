package yandex

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	timeofday "google.golang.org/genproto/googleapis/type/timeofday"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/postgresql/v1"
	config "github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/postgresql/v1/config"
)

type PostgreSQLHostSpec struct {
	HostSpec        *postgresql.HostSpec
	Fqdn            string
	HasComputedFqdn bool
}

func flattenPGClusterConfig(c *postgresql.ClusterConfig, d *schema.ResourceData) ([]interface{}, error) {
	poolerConf, err := flattenPGPoolerConfig(c.PoolerConfig)
	if err != nil {
		return nil, err
	}

	resources, err := flattenPGResources(c.Resources)
	if err != nil {
		return nil, err
	}

	backupWindowStart, err := flattenPGBackupWindowStart(c.BackupWindowStart)
	if err != nil {
		return nil, err
	}

	performanceDiagnostics, err := flattenPGPerformanceDiagnostics(c.PerformanceDiagnostics)
	if err != nil {
		return nil, err
	}

	access, err := flattenPGAccess(c.Access)
	if err != nil {
		return nil, err
	}

	settings, err := flattenPGSettings(c, d)
	if err != nil {
		return nil, err
	}

	out := map[string]interface{}{}
	out["autofailover"] = c.GetAutofailover().GetValue()
	out["version"] = c.Version
	out["pooler_config"] = poolerConf
	out["resources"] = resources
	out["backup_window_start"] = backupWindowStart
	out["performance_diagnostics"] = performanceDiagnostics
	out["access"] = access
	out["postgresql_config"] = settings

	return []interface{}{out}, nil
}

func flattenPGPoolerConfig(c *postgresql.ConnectionPoolerConfig) ([]interface{}, error) {
	if c == nil {
		return nil, nil
	}

	out := map[string]interface{}{}

	out["pool_discard"] = c.GetPoolDiscard().GetValue()
	out["pooling_mode"] = c.GetPoolingMode().String()

	return []interface{}{out}, nil
}

func flattenPGResources(r *postgresql.Resources) ([]interface{}, error) {
	out := map[string]interface{}{}
	out["resource_preset_id"] = r.ResourcePresetId
	out["disk_size"] = toGigabytes(r.DiskSize)
	out["disk_type_id"] = r.DiskTypeId

	return []interface{}{out}, nil
}

func flattenPGBackupWindowStart(t *timeofday.TimeOfDay) ([]interface{}, error) {
	if t == nil {
		return nil, nil
	}

	out := map[string]interface{}{}

	out["hours"] = int(t.Hours)
	out["minutes"] = int(t.Minutes)

	return []interface{}{out}, nil
}

func flattenPGPerformanceDiagnostics(p *postgresql.PerformanceDiagnostics) ([]interface{}, error) {
	if p == nil {
		return nil, nil
	}

	out := map[string]interface{}{}

	out["enabled"] = p.Enabled
	out["sessions_sampling_interval"] = int(p.SessionsSamplingInterval)
	out["statements_sampling_interval"] = int(p.StatementsSamplingInterval)

	return []interface{}{out}, nil
}

func flattenPGSettingsSPL(settings map[string]string, d *schema.ResourceData) map[string]string {
	spl, ok := d.GetOkExists("config.0.postgresql_config.shared_preload_libraries")
	if ok {
		if settings == nil {
			settings = make(map[string]string)
		}
		if _, exists := settings["shared_preload_libraries"]; !exists {
			settings["shared_preload_libraries"] = spl.(string)
		}
	}

	return settings
}

func flattenPGSettings(c *postgresql.ClusterConfig, d *schema.ResourceData) (map[string]string, error) {

	if cf, ok := c.PostgresqlConfig.(*postgresql.ClusterConfig_PostgresqlConfig_12); ok {
		settings, err := flattenResourceGenerateMapS(cf.PostgresqlConfig_12.UserConfig, false, mdbPGSettingsFieldsInfo, false, true)
		if err != nil {
			return nil, err
		}

		settings = flattenPGSettingsSPL(settings, d)

		return settings, err
	}
	if cf, ok := c.PostgresqlConfig.(*postgresql.ClusterConfig_PostgresqlConfig_12_1C); ok {
		settings, err := flattenResourceGenerateMapS(cf.PostgresqlConfig_12_1C.UserConfig, false, mdbPGSettingsFieldsInfo, false, true)
		if err != nil {
			return nil, err
		}

		settings = flattenPGSettingsSPL(settings, d)
		return settings, err
	}
	if cf, ok := c.PostgresqlConfig.(*postgresql.ClusterConfig_PostgresqlConfig_11); ok {
		settings, err := flattenResourceGenerateMapS(cf.PostgresqlConfig_11.UserConfig, false, mdbPGSettingsFieldsInfo, false, true)
		if err != nil {
			return nil, err
		}

		settings = flattenPGSettingsSPL(settings, d)
		return settings, err
	}
	if cf, ok := c.PostgresqlConfig.(*postgresql.ClusterConfig_PostgresqlConfig_11_1C); ok {
		settings, err := flattenResourceGenerateMapS(cf.PostgresqlConfig_11_1C.UserConfig, false, mdbPGSettingsFieldsInfo, false, true)
		if err != nil {
			return nil, err
		}

		settings = flattenPGSettingsSPL(settings, d)
		return settings, err
	}
	if cf, ok := c.PostgresqlConfig.(*postgresql.ClusterConfig_PostgresqlConfig_10); ok {
		settings, err := flattenResourceGenerateMapS(cf.PostgresqlConfig_10.UserConfig, false, mdbPGSettingsFieldsInfo, false, true)
		if err != nil {
			return nil, err
		}

		settings = flattenPGSettingsSPL(settings, d)
		return settings, err
	}
	if cf, ok := c.PostgresqlConfig.(*postgresql.ClusterConfig_PostgresqlConfig_10_1C); ok {
		settings, err := flattenResourceGenerateMapS(cf.PostgresqlConfig_10_1C.UserConfig, false, mdbPGSettingsFieldsInfo, false, true)
		if err != nil {
			return nil, err
		}

		settings = flattenPGSettingsSPL(settings, d)
		return settings, err
	}

	return nil, nil
}

func flattenPGAccess(a *postgresql.Access) ([]interface{}, error) {
	if a == nil {
		return nil, nil
	}

	out := map[string]interface{}{}

	out["data_lens"] = a.DataLens
	out["web_sql"] = a.WebSql

	return []interface{}{out}, nil
}

func flattenPGUsers(us []*postgresql.User, passwords map[string]string,
	fieldsInfo *objectFieldsInfo) ([]map[string]interface{}, error) {

	out := make([]map[string]interface{}, 0)

	for _, u := range us {
		ou, err := flattenPGUser(u, fieldsInfo)
		if err != nil {
			return nil, err
		}

		if v, ok := passwords[u.Name]; ok {
			ou["password"] = v
		}

		out = append(out, ou)
	}

	return out, nil
}

func flattenPGUser(u *postgresql.User,
	fieldsInfo *objectFieldsInfo) (map[string]interface{}, error) {

	m := map[string]interface{}{}
	m["name"] = u.Name
	m["login"] = u.GetLogin().GetValue()

	permissions, err := flattenPGUserPermissions(u.Permissions)
	if err != nil {
		return nil, err
	}
	m["permission"] = permissions

	m["grants"] = u.Grants

	m["conn_limit"] = u.ConnLimit

	settings, err := flattenResourceGenerateMapS(u.Settings, false, fieldsInfo, false, true)
	if err != nil {
		return nil, err
	}
	m["settings"] = settings

	return m, nil
}

func pgUsersPasswords(users []*postgresql.UserSpec) map[string]string {
	out := map[string]string{}
	for _, u := range users {
		out[u.Name] = u.Password
	}
	return out
}

func pgUserPermissionHash(v interface{}) int {
	m := v.(map[string]interface{})

	if n, ok := m["database_name"]; ok {
		return hashcode.String(n.(string))
	}
	return 0
}

func flattenPGUserPermissions(ps []*postgresql.Permission) (*schema.Set, error) {
	out := schema.NewSet(pgUserPermissionHash, nil)

	for _, p := range ps {
		op := map[string]interface{}{
			"database_name": p.DatabaseName,
		}

		out.Add(op)
	}

	return out, nil
}

func flattenPGHosts(hs []*postgresql.Host) ([]map[string]interface{}, error) {
	out := []map[string]interface{}{}

	for _, h := range hs {
		m := map[string]interface{}{}
		m["zone"] = h.ZoneId
		m["subnet_id"] = h.SubnetId
		m["assign_public_ip"] = h.AssignPublicIp
		m["fqdn"] = h.Name
		m["role"] = h.Role.String()

		out = append(out, m)
	}

	return out, nil
}

func flattenPGDatabases(dbs []*postgresql.Database) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)

	for _, d := range dbs {
		m := make(map[string]interface{})
		m["name"] = d.Name
		m["owner"] = d.Owner
		m["lc_collate"] = d.LcCollate
		m["lc_type"] = d.LcCtype
		m["extension"] = flattenPGExtensions(d.Extensions)

		out = append(out, m)
	}

	return out
}

func flattenPGExtensions(es []*postgresql.Extension) *schema.Set {
	out := schema.NewSet(pgExtensionHash, nil)

	for _, e := range es {
		m := make(map[string]interface{})
		m["name"] = e.Name
		m["version"] = e.Version

		out.Add(m)
	}

	return out
}

func pgExtensionHash(v interface{}) int {
	var buf bytes.Buffer

	m := v.(map[string]interface{})
	if v, ok := m["name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["version"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func expandPGConfigSpec(d *schema.ResourceData) (cs *postgresql.ConfigSpec, updateFieldConfigName string, err error) {
	cs = &postgresql.ConfigSpec{}

	if v, ok := d.GetOk("config.0.version"); ok {
		cs.Version = v.(string)
	}

	if v, ok := d.GetOkExists("config.0.autofailover"); ok {
		cs.Autofailover = &wrappers.BoolValue{Value: v.(bool)}
	}

	poolerConfig, err := expandPGPoolerConfig(d)
	if err != nil {
		return nil, updateFieldConfigName, err
	}
	cs.PoolerConfig = poolerConfig

	resources, err := expandPGResources(d)
	if err != nil {
		return nil, updateFieldConfigName, err
	}
	cs.Resources = resources

	cs.BackupWindowStart = expandPGBackupWindowStart(d)
	cs.Access = expandPGAccess(d)
	cs.PerformanceDiagnostics = expandPGPerformanceDiagnostics(d)

	updateFieldConfigName, err = expandPGConfigSpecSettings(d, cs)
	if err != nil {
		return nil, updateFieldConfigName, err
	}

	return cs, updateFieldConfigName, nil
}

func expandPGPoolerConfig(d *schema.ResourceData) (*postgresql.ConnectionPoolerConfig, error) {
	pc := &postgresql.ConnectionPoolerConfig{}

	if v, ok := d.GetOk("config.0.pooler_config.0.pooling_mode"); ok {
		pm, err := parsePostgreSQLPoolingMode(v.(string))
		if err != nil {
			return nil, err
		}

		pc.PoolingMode = pm
	}

	if v, ok := d.GetOk("config.0.pooler_config.0.pool_discard"); ok {
		pc.PoolDiscard = &wrappers.BoolValue{Value: v.(bool)}
	}

	return pc, nil
}

func expandPGResources(d *schema.ResourceData) (*postgresql.Resources, error) {
	r := &postgresql.Resources{}

	if v, ok := d.GetOk("config.0.resources.0.resource_preset_id"); ok {
		r.ResourcePresetId = v.(string)
	}

	if v, ok := d.GetOk("config.0.resources.0.disk_size"); ok {
		r.DiskSize = toBytes(v.(int))
	}

	if v, ok := d.GetOk("config.0.resources.0.disk_type_id"); ok {
		r.DiskTypeId = v.(string)
	}

	return r, nil
}

func expandPGUserSpecs(d *schema.ResourceData) ([]*postgresql.UserSpec, error) {
	out := []*postgresql.UserSpec{}

	cnt := d.Get("user.#").(int)
	for i := 0; i < cnt; i++ {
		user, err := expandPGUserNew(d, fmt.Sprintf("user.%v.", i))
		if err != nil {
			return nil, err
		}

		out = append(out, user)
	}

	return out, nil
}

// pgUserForCreate get users for create
func pgUserForCreate(d *schema.ResourceData, currUsers []*postgresql.User) (usersForCreate []*postgresql.UserSpec, err error) {
	currentUser := make(map[string]struct{})
	for _, v := range currUsers {
		currentUser[v.Name] = struct{}{}
	}
	usersForCreate = make([]*postgresql.UserSpec, 0)

	cnt := d.Get("user.#").(int)
	for i := 0; i < cnt; i++ {
		_, ok := currentUser[d.Get(fmt.Sprintf("user.%v.name", i)).(string)]
		if !ok {
			user, err := expandPGUserNew(d, fmt.Sprintf("user.%v.", i))
			if err != nil {
				return nil, err
			}
			user.Grants = make([]string, 0)
			user.Permissions = make([]*postgresql.Permission, 0)
			usersForCreate = append(usersForCreate, user)
		}
	}

	return usersForCreate, nil
}

// expandPGUserNew expand to new object from schema.ResourceData
// path like "user.3."
func expandPGUserNew(d *schema.ResourceData, path string) (*postgresql.UserSpec, error) {
	return expandPGUser(d, &postgresql.UserSpec{}, path)
}

// expandPGUser expand to exists object from schema.ResourceData
// path like "user.3."
func expandPGUser(d *schema.ResourceData, user *postgresql.UserSpec, path string) (*postgresql.UserSpec, error) {

	if v, ok := d.GetOkExists(path + "name"); ok {
		user.Name = v.(string)
	}

	if v, ok := d.GetOkExists(path + "password"); ok {
		user.Password = v.(string)
	}

	if v, ok := d.GetOkExists(path + "login"); ok {
		user.Login = &wrappers.BoolValue{Value: v.(bool)}
	}

	if v, ok := d.GetOkExists(path + "conn_limit"); ok {
		user.ConnLimit = &wrappers.Int64Value{Value: int64(v.(int))}
	}

	if v, ok := d.GetOkExists(path + "permission"); ok {
		permissions, err := expandPGUserPermissions(v.(*schema.Set))
		if err != nil {
			return nil, err
		}
		user.Permissions = permissions
	}

	if v, ok := d.GetOkExists(path + "grants"); ok {
		gs, err := expandPGUserGrants(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		user.Grants = gs
	}

	if _, ok := d.GetOkExists(path + "settings"); ok {
		if user.Settings == nil {
			user.Settings = &postgresql.UserSettings{}
		}

		err := expandResourceGenerate(mdbPGUserSettingsFieldsInfo, d, user.Settings, path+"settings.", true)
		if err != nil {
			return nil, err
		}

	}

	return user, nil
}

func expandPGUserGrants(gs []interface{}) ([]string, error) {
	out := []string{}

	if gs == nil {
		return out, nil
	}

	for _, v := range gs {
		out = append(out, v.(string))
	}

	return out, nil
}

func expandPGUserPermissions(ps *schema.Set) ([]*postgresql.Permission, error) {
	out := []*postgresql.Permission{}

	for _, p := range ps.List() {
		m := p.(map[string]interface{})
		permission := &postgresql.Permission{}

		if v, ok := m["database_name"]; ok {
			permission.DatabaseName = v.(string)
		}

		out = append(out, permission)
	}

	return out, nil
}

func expandPGHosts(d *schema.ResourceData) ([]*PostgreSQLHostSpec, error) {
	out := []*PostgreSQLHostSpec{}
	hosts := d.Get("host").([]interface{})

	for _, v := range hosts {
		m := v.(map[string]interface{})
		h, err := expandPGHost(m)
		if err != nil {
			return nil, err
		}
		out = append(out, h)
	}

	return out, nil
}

func expandPGHost(m map[string]interface{}) (*PostgreSQLHostSpec, error) {
	hostSpec := &postgresql.HostSpec{}
	host := &PostgreSQLHostSpec{HostSpec: hostSpec, HasComputedFqdn: false}
	if v, ok := m["zone"]; ok {
		host.HostSpec.ZoneId = v.(string)
	}

	if v, ok := m["subnet_id"]; ok {
		host.HostSpec.SubnetId = v.(string)
	}

	if v, ok := m["assign_public_ip"]; ok {
		host.HostSpec.AssignPublicIp = v.(bool)
	}
	if v, ok := m["fqdn"]; ok && v.(string) != "" {
		host.HasComputedFqdn = true
		host.Fqdn = v.(string)
	}

	return host, nil
}

func sortPGHosts(hosts []*postgresql.Host, specs []*PostgreSQLHostSpec) {
	for i, h := range specs {
		for j := i + 1; j < len(hosts); j++ {
			if h.HostSpec.ZoneId == hosts[j].ZoneId {
				hosts[i], hosts[j] = hosts[j], hosts[i]
				break
			}
		}
	}
}

func expandPGDatabaseSpecs(d *schema.ResourceData) ([]*postgresql.DatabaseSpec, error) {
	out := []*postgresql.DatabaseSpec{}
	dbs := d.Get("database").([]interface{})

	for _, d := range dbs {
		m := d.(map[string]interface{})
		database, err := expandPGDatabase(m)
		if err != nil {
			return nil, err
		}

		out = append(out, database)
	}

	return out, nil
}

func expandPGDatabase(m map[string]interface{}) (*postgresql.DatabaseSpec, error) {
	out := &postgresql.DatabaseSpec{}

	if v, ok := m["name"]; ok {
		out.Name = v.(string)
	}

	if v, ok := m["owner"]; ok {
		out.Owner = v.(string)
	}

	if v, ok := m["lc_collate"]; ok {
		out.LcCollate = v.(string)
	}

	if v, ok := m["lc_type"]; ok {
		out.LcCtype = v.(string)
	}

	if v, ok := m["extension"]; ok {
		es := v.(*schema.Set).List()
		extensions, err := expandPGExtensions(es)
		if err != nil {
			return nil, err
		}

		out.Extensions = extensions
	}

	return out, nil
}

func expandPGExtensions(es []interface{}) ([]*postgresql.Extension, error) {
	out := []*postgresql.Extension{}

	for _, e := range es {
		m := e.(map[string]interface{})
		extension := &postgresql.Extension{}

		if v, ok := m["name"]; ok {
			extension.Name = v.(string)
		}

		if v, ok := m["version"]; ok {
			extension.Version = v.(string)
		}

		out = append(out, extension)
	}

	return out, nil
}

func expandPGBackupWindowStart(d *schema.ResourceData) *timeofday.TimeOfDay {
	out := &timeofday.TimeOfDay{}

	if v, ok := d.GetOk("config.0.backup_window_start.0.hours"); ok {
		out.Hours = int32(v.(int))
	}

	if v, ok := d.GetOk("config.0.backup_window_start.0.minutes"); ok {
		out.Minutes = int32(v.(int))
	}

	return out
}

func expandPGPerformanceDiagnostics(d *schema.ResourceData) *postgresql.PerformanceDiagnostics {

	if _, ok := d.GetOkExists("config.0.performance_diagnostics"); !ok {
		return nil
	}

	out := &postgresql.PerformanceDiagnostics{}

	if v, ok := d.GetOk("config.0.performance_diagnostics.0.enabled"); ok {
		out.Enabled = v.(bool)
	}

	if v, ok := d.GetOk("config.0.performance_diagnostics.0.sessions_sampling_interval"); ok {
		out.SessionsSamplingInterval = int64(v.(int))
	}

	if v, ok := d.GetOk("config.0.performance_diagnostics.0.statements_sampling_interval"); ok {
		out.StatementsSamplingInterval = int64(v.(int))
	}

	return out
}

func expandPGAccess(d *schema.ResourceData) *postgresql.Access {
	out := &postgresql.Access{}

	if v, ok := d.GetOk("config.0.access.0.data_lens"); ok {
		out.DataLens = v.(bool)
	}

	if v, ok := d.GetOk("config.0.access.0.web_sql"); ok {
		out.WebSql = v.(bool)
	}

	return out
}

func expandPGConfigSpecSettings(d *schema.ResourceData, configSpec *postgresql.ConfigSpec) (updateFieldConfigName string, err error) {

	version := configSpec.Version

	path := "config.0.postgresql_config"

	if _, ok := d.GetOkExists(path); ok {

		var sharedPreloadLibraries []int32
		sharedPreloadLibValue, ok := d.GetOkExists(path + ".shared_preload_libraries")
		if ok {
			splValue := sharedPreloadLibValue.(string)

			for _, sv := range strings.Split(splValue, ",") {

				i, err := mdbPGSettingsFieldsInfo.stringToInt("shared_preload_libraries", sv)
				if err != nil {
					return updateFieldConfigName, err
				}
				if i != nil {
					sharedPreloadLibraries = append(sharedPreloadLibraries, int32(*i))
				}
			}
		}

		var pgConfig interface{}
		if version == "10" {
			cfg := &postgresql.ConfigSpec_PostgresqlConfig_10{
				PostgresqlConfig_10: &config.PostgresqlConfig10{},
			}
			if len(sharedPreloadLibraries) > 0 {
				for _, v := range sharedPreloadLibraries {
					cfg.PostgresqlConfig_10.SharedPreloadLibraries = append(cfg.PostgresqlConfig_10.SharedPreloadLibraries, config.PostgresqlConfig10_SharedPreloadLibraries(v))
				}
			}
			pgConfig = cfg.PostgresqlConfig_10
			configSpec.PostgresqlConfig = cfg
			updateFieldConfigName = "postgresql_config_10"
		} else if version == "10-1c" {
			cfg := &postgresql.ConfigSpec_PostgresqlConfig_10_1C{
				PostgresqlConfig_10_1C: &config.PostgresqlConfig10_1C{},
			}
			if len(sharedPreloadLibraries) > 0 {
				for _, v := range sharedPreloadLibraries {
					cfg.PostgresqlConfig_10_1C.SharedPreloadLibraries = append(cfg.PostgresqlConfig_10_1C.SharedPreloadLibraries, config.PostgresqlConfig10_1C_SharedPreloadLibraries(v))
				}
			}
			pgConfig = cfg.PostgresqlConfig_10_1C
			configSpec.PostgresqlConfig = cfg
			updateFieldConfigName = "postgresql_config_10_1c"
		} else if version == "11" {
			cfg := &postgresql.ConfigSpec_PostgresqlConfig_11{
				PostgresqlConfig_11: &config.PostgresqlConfig11{},
			}
			if len(sharedPreloadLibraries) > 0 {
				for _, v := range sharedPreloadLibraries {
					cfg.PostgresqlConfig_11.SharedPreloadLibraries = append(cfg.PostgresqlConfig_11.SharedPreloadLibraries, config.PostgresqlConfig11_SharedPreloadLibraries(v))
				}
			}
			pgConfig = cfg.PostgresqlConfig_11
			configSpec.PostgresqlConfig = cfg
			updateFieldConfigName = "postgresql_config_11"
		} else if version == "11-1c" {
			cfg := &postgresql.ConfigSpec_PostgresqlConfig_11_1C{
				PostgresqlConfig_11_1C: &config.PostgresqlConfig11_1C{},
			}
			if len(sharedPreloadLibraries) > 0 {
				for _, v := range sharedPreloadLibraries {
					cfg.PostgresqlConfig_11_1C.SharedPreloadLibraries = append(cfg.PostgresqlConfig_11_1C.SharedPreloadLibraries, config.PostgresqlConfig11_1C_SharedPreloadLibraries(v))
				}
			}
			pgConfig = cfg.PostgresqlConfig_11_1C
			configSpec.PostgresqlConfig = cfg
			updateFieldConfigName = "postgresql_config_11_1c"
		} else if version == "12-1c" {
			cfg := &postgresql.ConfigSpec_PostgresqlConfig_12_1C{
				PostgresqlConfig_12_1C: &config.PostgresqlConfig12_1C{},
			}
			if len(sharedPreloadLibraries) > 0 {
				for _, v := range sharedPreloadLibraries {
					cfg.PostgresqlConfig_12_1C.SharedPreloadLibraries = append(cfg.PostgresqlConfig_12_1C.SharedPreloadLibraries, config.PostgresqlConfig12_1C_SharedPreloadLibraries(v))
				}
			}
			pgConfig = cfg.PostgresqlConfig_12_1C
			configSpec.PostgresqlConfig = cfg
			updateFieldConfigName = "postgresql_config_12_1c"
		} else {
			// 12
			cfg := &postgresql.ConfigSpec_PostgresqlConfig_12{
				PostgresqlConfig_12: &config.PostgresqlConfig12{},
			}
			if len(sharedPreloadLibraries) > 0 {
				for _, v := range sharedPreloadLibraries {
					cfg.PostgresqlConfig_12.SharedPreloadLibraries = append(cfg.PostgresqlConfig_12.SharedPreloadLibraries, config.PostgresqlConfig12_SharedPreloadLibraries(v))
				}
			}
			pgConfig = cfg.PostgresqlConfig_12
			configSpec.PostgresqlConfig = cfg
			updateFieldConfigName = "postgresql_config_12"
		}

		err := expandResourceGenerate(mdbPGSettingsFieldsInfo, d, pgConfig, path+".", true)

		if err != nil {
			return updateFieldConfigName, err
		}

	}

	return updateFieldConfigName, nil
}

func pgDatabasesDiff(currDBs []*postgresql.Database, targetDBs []*postgresql.DatabaseSpec) ([]string, []*postgresql.DatabaseSpec) {
	m := map[string]bool{}
	toAdd := []*postgresql.DatabaseSpec{}
	toDelete := map[string]bool{}
	for _, db := range currDBs {
		toDelete[db.Name] = true
		m[db.Name] = true
	}

	for _, db := range targetDBs {
		delete(toDelete, db.Name)
		if _, ok := m[db.Name]; !ok {
			toAdd = append(toAdd, db)
		}
	}

	toDel := []string{}
	for u := range toDelete {
		toDel = append(toDel, u)
	}

	return toDel, toAdd
}

func pgChangedDatabases(oldSpecs []interface{}, newSpecs []interface{}) ([]*postgresql.DatabaseSpec, error) {
	out := []*postgresql.DatabaseSpec{}

	m := map[string]*postgresql.DatabaseSpec{}
	for _, spec := range oldSpecs {
		db, err := expandPGDatabase(spec.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		m[db.Name] = db
	}

	for _, spec := range newSpecs {
		db, err := expandPGDatabase(spec.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		if oldDB, ok := m[db.Name]; ok {
			if !reflect.DeepEqual(db, oldDB) {
				out = append(out, db)
			}
		}
	}

	return out, nil
}

func pgHostsDiff(currHosts []*postgresql.Host, targetHosts []*PostgreSQLHostSpec) ([]string, []*postgresql.HostSpec) {
	m := map[string]*PostgreSQLHostSpec{}

	toAdd := []*postgresql.HostSpec{}
	for _, h := range targetHosts {
		if !h.HasComputedFqdn {
			toAdd = append(toAdd, h.HostSpec)
		} else {
			m[h.Fqdn] = h
		}
	}

	toDelete := []string{}
	for _, h := range currHosts {
		_, ok := m[h.Name]
		if !ok {
			toDelete = append(toDelete, h.Name)
		}
	}

	return toDelete, toAdd
}

func parsePostgreSQLEnv(e string) (postgresql.Cluster_Environment, error) {
	v, ok := postgresql.Cluster_Environment_value[e]
	if !ok {
		return 0, fmt.Errorf("value for 'environment' must be one of %s, not `%s`",
			getJoinedKeys(getEnumValueMapKeys(postgresql.Cluster_Environment_value)), e)
	}

	return postgresql.Cluster_Environment(v), nil
}

func parsePostgreSQLPoolingMode(s string) (postgresql.ConnectionPoolerConfig_PoolingMode, error) {
	v, ok := postgresql.ConnectionPoolerConfig_PoolingMode_value[s]
	if !ok {
		return 0, fmt.Errorf("value for 'pooling_mode' must be one of %s, not `%s`",
			getJoinedKeys(getEnumValueMapKeys(postgresql.ConnectionPoolerConfig_PoolingMode_value)), s)
	}

	return postgresql.ConnectionPoolerConfig_PoolingMode(v), nil
}

func mdbPGSharedPreloadLibrariesCheck(fieldsInfo *objectFieldsInfo, v interface{}) error {

	s, ok := v.(string)
	if ok {
		if s == "" {
			return nil
		}

		for _, sv := range strings.Split(s, ",") {

			i, err := fieldsInfo.stringToInt("shared_preload_libraries", sv)
			if err != nil {
				return err
			}
			err = fieldsInfo.intCheckSetValue("shared_preload_libraries", i)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("mdbPGSharedPreloadLibrariesCheck: Unsupported type for value %v", v)
}
func mdbPGSharedPreloadLibrariesCompare(fieldsInfo *objectFieldsInfo, old, new string) bool {
	return old == new
}

var mdbPGUserSettingsTransactionIsolationName = map[int]string{
	int(postgresql.UserSettings_TRANSACTION_ISOLATION_UNSPECIFIED):      "unspecified",
	int(postgresql.UserSettings_TRANSACTION_ISOLATION_READ_UNCOMMITTED): "read uncommitted",
	int(postgresql.UserSettings_TRANSACTION_ISOLATION_READ_COMMITTED):   "read committed",
	int(postgresql.UserSettings_TRANSACTION_ISOLATION_REPEATABLE_READ):  "repeatable read",
	int(postgresql.UserSettings_TRANSACTION_ISOLATION_SERIALIZABLE):     "serializable",
}
var mdbPGUserSettingsSynchronousCommitName = map[int]string{
	int(postgresql.UserSettings_SYNCHRONOUS_COMMIT_UNSPECIFIED):  "unspecified",
	int(postgresql.UserSettings_SYNCHRONOUS_COMMIT_ON):           "on",
	int(postgresql.UserSettings_SYNCHRONOUS_COMMIT_OFF):          "off",
	int(postgresql.UserSettings_SYNCHRONOUS_COMMIT_LOCAL):        "local",
	int(postgresql.UserSettings_SYNCHRONOUS_COMMIT_REMOTE_WRITE): "remote write",
	int(postgresql.UserSettings_SYNCHRONOUS_COMMIT_REMOTE_APPLY): "remote apply",
}
var mdbPGUserSettingsLogStatementName = map[int]string{
	int(postgresql.UserSettings_LOG_STATEMENT_UNSPECIFIED): "unspecified",
	int(postgresql.UserSettings_LOG_STATEMENT_NONE):        "none",
	int(postgresql.UserSettings_LOG_STATEMENT_DDL):         "ddl",
	int(postgresql.UserSettings_LOG_STATEMENT_MOD):         "mod",
	int(postgresql.UserSettings_LOG_STATEMENT_ALL):         "all",
}

var mdbPGUserSettingsFieldsInfo = newObjectFieldsInfo().
	addType(postgresql.UserSettings{}).
	addIDefault("log_min_duration_statement", -1).
	addEnumHumanNames("default_transaction_isolation", mdbPGUserSettingsTransactionIsolationName,
		postgresql.UserSettings_TransactionIsolation_name).
	addEnumHumanNames("synchronous_commit", mdbPGUserSettingsSynchronousCommitName,
		postgresql.UserSettings_SynchronousCommit_name).
	addEnumHumanNames("log_statement", mdbPGUserSettingsLogStatementName,
		postgresql.UserSettings_LogStatement_name)

var mdbPGSettingsFieldsInfo = newObjectFieldsInfo().
	addType(config.PostgresqlConfig12{}).
	addType(config.PostgresqlConfig12_1C{}).
	addType(config.PostgresqlConfig11{}).
	addType(config.PostgresqlConfig11_1C{}).
	addType(config.PostgresqlConfig10{}).
	addType(config.PostgresqlConfig10_1C{}).
	addEnumGeneratedNames("wal_level", config.PostgresqlConfig12_WalLevel_name).
	addEnumGeneratedNames("synchronous_commit", config.PostgresqlConfig12_SynchronousCommit_name).
	addEnumGeneratedNames("constraint_exclusion", config.PostgresqlConfig12_ConstraintExclusion_name).
	addEnumGeneratedNames("force_parallel_mode", config.PostgresqlConfig12_ForceParallelMode_name).
	addEnumGeneratedNames("client_min_messages", config.PostgresqlConfig12_LogLevel_name).
	addEnumGeneratedNames("log_min_messages", config.PostgresqlConfig12_LogLevel_name).
	addEnumGeneratedNames("log_min_error_statement", config.PostgresqlConfig12_LogLevel_name).
	addEnumGeneratedNames("log_error_verbosity", config.PostgresqlConfig12_LogErrorVerbosity_name).
	addEnumGeneratedNames("log_statement", config.PostgresqlConfig12_LogStatement_name).
	addEnumGeneratedNames("default_transaction_isolation", config.PostgresqlConfig12_TransactionIsolation_name).
	addEnumGeneratedNames("bytea_output", config.PostgresqlConfig12_ByteaOutput_name).
	addEnumGeneratedNames("xmlbinary", config.PostgresqlConfig12_XmlBinary_name).
	addEnumGeneratedNames("xmloption", config.PostgresqlConfig12_XmlOption_name).
	addEnumGeneratedNames("backslash_quote", config.PostgresqlConfig12_BackslashQuote_name).
	addEnumGeneratedNames("plan_cache_mode", config.PostgresqlConfig12_PlanCacheMode_name).
	addSkipEnumGeneratedNames("shared_preload_libraries", config.PostgresqlConfig12_SharedPreloadLibraries_name, mdbPGSharedPreloadLibrariesCheck, mdbPGSharedPreloadLibrariesCompare).
	addEnumGeneratedNames("pg_hint_plan_debug_print", config.PostgresqlConfig12_PgHintPlanDebugPrint_name).
	addEnumGeneratedNames("pg_hint_plan_message_level", config.PostgresqlConfig12_LogLevel_name)
