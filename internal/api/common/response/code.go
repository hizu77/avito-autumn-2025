//go:generate go-enum  --output-suffix=.generated

package response

// ENUM(
//
//		TeamExists=TEAM_EXISTS,
//		AdminExists=ADMIN_EXISTS
//		PrExists=PR_EXISTS,
//		PrMerged=PR_MERGED,
//		NotAssigned=NOT_ASSIGNED,
//		NoCandidate=NO_CANDIDATE,
//		NotFound=NOT_FOUND,
//	 	Internal=INTERNAL,
//		BadRequest=BAD_REQUEST,
//		InvalidCredentials=INVALID_CREDS
//		Unauthorized=UNAUTHORIZED,
//
// )
type Code string
