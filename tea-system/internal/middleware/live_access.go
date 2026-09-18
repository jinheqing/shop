package middleware

import (
	"net/http"
	"strconv"
	"tea-system/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LiveAccessMiddleware — 检查当前用户是否能访问某个直播
func LiveAccessMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": 400, "message": "missing id"})
			return
		}
		var room models.LiveRoom
		if err := db.First(&room, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": 404, "message": "live room not found"})
			return
		}

		switch room.Visibility {
		case "public":
			c.Set("live_room", &room)
			c.Next()
			return
		case "registered":
			if _, ok := c.Get("user_id"); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code": 401, "message": "login required to access this broadcast", "visibility": room.Visibility,
				})
				return
			}
			c.Set("live_room", &room)
			c.Next()
			return
		case "restricted":
			userIDRaw, ok := c.Get("user_id")
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "this broadcast is private"})
				return
			}
			userID, _ := userIDRaw.(uint64)
			if roomVisibleToUser(db, &room, userID) {
				c.Set("live_room", &room)
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 403, "message": "this broadcast is private — your advisor has not granted access",
			})
			return
		}
		c.Set("live_room", &room)
		c.Next()
	}
}

// RecordingAccessMiddleware — 同上，针对 Recording
func RecordingAccessMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var rec models.Recording
		if err := db.First(&rec, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": 404, "message": "recording not found"})
			return
		}
		switch rec.Visibility {
		case "public":
			c.Set("recording", &rec)
			c.Next()
			return
		case "registered":
			if _, ok := c.Get("user_id"); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "login required"})
				return
			}
			c.Set("recording", &rec)
			c.Next()
			return
		case "restricted":
			userIDRaw, ok := c.Get("user_id")
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "private recording"})
				return
			}
			userID, _ := userIDRaw.(uint64)
			if recordingVisibleToUser(db, &rec, userID) {
				c.Set("recording", &rec)
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "private recording"})
			return
		}
		c.Set("recording", &rec)
		c.Next()
	}
}

// === 辅助函数 ===

func roomVisibleToUser(db *gorm.DB, room *models.LiveRoom, userID uint64) bool {
	for _, uid := range room.VisibleUserIDs {
		if v, err := strconv.ParseUint(uid, 10, 64); err == nil && v == userID {
			return true
		}
	}
	return userInAnyGroup(db, userID, room.VisibleGroupIDs)
}

func recordingVisibleToUser(db *gorm.DB, rec *models.Recording, userID uint64) bool {
	for _, uid := range rec.VisibleUserIDs {
		if v, err := strconv.ParseUint(uid, 10, 64); err == nil && v == userID {
			return true
		}
	}
	return userInAnyGroup(db, userID, rec.VisibleGroupIDs)
}

func userInAnyGroup(db *gorm.DB, userID uint64, groupIDs []string) bool {
	if len(groupIDs) == 0 {
		return false
	}
	var userGroups []models.UserGroupMember
	db.Where("user_id = ?", userID).Find(&userGroups)
	set := map[uint64]bool{}
	for _, m := range userGroups {
		set[m.GroupID] = true
	}
	for _, gid := range groupIDs {
		if v, err := strconv.ParseUint(gid, 10, 64); err == nil && set[v] {
			return true
		}
	}
	return false
}
