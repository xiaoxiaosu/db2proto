package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"xiaoxiaosu.com/db2pb/internal/db2pb"
)

var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string

var pbCmd = &cobra.Command{
	Use: "pb",
	Short: "表结构转为pb文件",
	Long: "表结构转为pb文件",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo := &db2pb.DBInfo{
			DBType: dbType,
			Host: host,
			UserName: username,
			Password: password,
			Charset: charset,
		}
		dbModel := db2pb.NewDBModel(dbInfo)
		err := dbModel.Connect()
		if err != nil {
			log.Fatalf("dbModel.Connect err:%v", err)
		}
		columns, err := dbModel.GetColumns(dbName, tableName)
		if err != nil {
			log.Fatalf("dbModel.GetColumns err:%v", err)
		}

		template := db2pb.NewStructTemplate()
		templateColumns := template.AssemblyColumns(columns)
		err = template.Generate(tableName, templateColumns)
		if err != nil {
			log.Fatalf("template.Generate err:%v", err)
		}
	},
}



func init() {
	pbCmd.Flags().StringVarP(&username, "username", "", "root", "请输入数据库username")
	pbCmd.Flags().StringVarP(&password, "password", "", "", "请输入数据库password")
	pbCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1", "请输入数据库HOST")
	pbCmd.Flags().StringVarP(&charset, "charset", "", "utf8", "请输入数据库HOST")
	pbCmd.Flags().StringVarP(&dbType, "type", "", "mysql", "请输入数据库类型")
	pbCmd.Flags().StringVarP(&dbName, "db", "", "mysql", "请输入数据库名称")
	pbCmd.Flags().StringVarP(&tableName, "table", "", "mysql", "请输入表名")
}