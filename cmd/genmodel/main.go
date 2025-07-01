package main

import (
	"financing-offer/assets"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"

	"financing-offer/internal/config"
)

const defaultSchema = "public"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg, err := config.InitConfig[config.AppConfig](assets.EmbeddedFiles)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
		return
	}
	port, err := strconv.Atoi(cfg.Db.Port)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
		return
	}
	dbConn := postgres.DBConnection{
		Host:       cfg.Db.Host,
		Port:       port,
		User:       cfg.Db.User,
		Password:   cfg.Db.Password,
		SslMode:    "disable",
		DBName:     cfg.Db.DbName,
		SchemaName: defaultSchema,
	}

	err = postgres.Generate(
		cfg.ModelGeneration.Path,
		dbConn,
		template.Default(postgres2.Dialect).
			UseSchema(
				func(schema metadata.Schema) template.Schema {
					return template.DefaultSchema(schema).
						UseModel(
							template.DefaultModel().
								UseTable(
									func(table metadata.Table) template.TableModel {
										tableModel := template.DefaultTableModel(table)
										// skip schema_migrations table
										if schema.Name == defaultSchema && slices.Contains(
											cfg.ModelGeneration.IgnoredTables, table.Name,
										) {
											tableModel.Skip = true
										}
										return tableModel.UseField(
											func(column metadata.Column) template.TableModelField {
												defaultTableModelField := template.DefaultTableModelField(column)
												if column.DataType.Name == "numeric" {
													defaultTableModelField.Type = template.NewType(decimal.Decimal{})
												}
												if column.DataType.Name == "timestamp without time zone" && column.IsNullable {
													defaultTableModelField.Type = template.NewType(null.Time{})
												}
												return defaultTableModelField
											},
										)
									},
								),
						).UseSQLBuilder(
						template.DefaultSQLBuilder().
							UseTable(
								func(table metadata.Table) template.TableSQLBuilder {
									sqlBuilder := template.DefaultTableSQLBuilder(table)
									if schema.Name == defaultSchema && slices.Contains(
										cfg.ModelGeneration.IgnoredTables, table.Name,
									) {
										sqlBuilder.Skip = true
									}
									return sqlBuilder
								},
							),
					)
				},
			),
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
		return
	}
	pathToTables := fmt.Sprintf("%s/%s/%s/table", cfg.ModelGeneration.Path, cfg.Db.DbName, defaultSchema)
	if err := AutoImmutableColumns(pathToTables, cfg.ModelGeneration.ImmutableColumns); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
