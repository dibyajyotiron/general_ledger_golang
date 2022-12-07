package rpc

import (
	"context"
	"encoding/json"
	"errors"

	proto "general_ledger_golang/api/proto/code/go"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/util"

	"general_ledger_golang/service/book_service"
)

// GetBalance returns a map where key is the asset name, and value is the amount of that asset in that book
func (*Grpc) GetBalance(_ context.Context, req *proto.GetBalanceRequest) (res *proto.GetBalanceResponse, err error) {
	logger.Logger.Infof("Invoked GetBalance")
	bookService := book_service.BookService{}

	result, err := bookService.GetBalance(req.BookId, "", "", nil)
	marshal, _ := json.Marshal(result)

	logger.Logger.Infof("Result: %+v", string(marshal))

	if err != nil {
		return nil, err
	}

	resBalances := map[string]string{}
	for k, v := range result {
		vMap, err := util.InterfaceToMapOfString(v)
		// in case of err, return
		if err != nil {
			logger.Logger.Infof("Error Occured while converting interface to map of string, v: %+v, err: %+v", v, err)
			return nil, errors.New("something went wrong, we're checking")
		}
		resBalances[k] = vMap["balance"]
	}
	return &proto.GetBalanceResponse{
		Balances: resBalances,
	}, nil
}
