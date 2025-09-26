package interceptors

import pb "soa-socialnetwork/services/accounts/proto"

type authRequirements struct {
	needAuth        bool
	needReadAccess  bool
	needWriteAccess bool
}

var methods_auth_requirements = map[string]authRequirements{
	pb.AccountsService_RegisterUser_FullMethodName: {
		needAuth:        false,
		needReadAccess:  false,
		needWriteAccess: false,
	},
	pb.AccountsService_UnregisterUser_FullMethodName: {
		needAuth:        true,
		needReadAccess:  false,
		needWriteAccess: true,
	},
	pb.AccountsService_GetProfile_FullMethodName: {
		needAuth:        false,
		needReadAccess:  false,
		needWriteAccess: false,
	},
	pb.AccountsService_EditProfile_FullMethodName: {
		needAuth:        true,
		needReadAccess:  false,
		needWriteAccess: true,
	},
	pb.AccountsService_Authenticate_FullMethodName: {
		needAuth:        false,
		needReadAccess:  false,
		needWriteAccess: false,
	},
	pb.AccountsService_CreateApiToken_FullMethodName: {
		needAuth:        false,
		needReadAccess:  false,
		needWriteAccess: false,
	},
	pb.AccountsService_ValidateApiToken_FullMethodName: {
		needAuth:        false,
		needReadAccess:  false,
		needWriteAccess: false,
	},
	pb.AccountsService_ResolveProfileId_FullMethodName: {
		needAuth:        false,
		needReadAccess:  false,
		needWriteAccess: false,
	},
	pb.AccountsService_ResolveAccountId_FullMethodName: {
		needAuth:        false,
		needReadAccess:  false,
		needWriteAccess: false,
	},
}

func getAuthRequirements(fullMethodName string) authRequirements {
	reqs, ok := methods_auth_requirements[fullMethodName]
	if !ok {
		return authRequirements{
			needAuth:        false,
			needReadAccess:  false,
			needWriteAccess: false,
		}
	}

	return reqs
}
