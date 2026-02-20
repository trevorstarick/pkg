package mediainfo

import "slices"

var KeyGroupCommon = slices.Concat(
	KeyGroupCodec,
	KeyGroupDates,
	KeyGroupFormat,
	[]string{
		KeyExtra,
		KeyID,
		KeyLanguage,
		KeyType,
		"@typeorder",
		"Title",
		"UniqueID",
	},
)

var KeyGroupGeneral = slices.Concat(
	KeyGroupCommon,
	KeyGroupOverallBitRate,
	KeyGroupFrame,
	KeyGroupTrackCount,
	KeyGroupStream,
	KeyGroupEncoded,
	KeyGroupTags,
	[]string{
		KeyDataSize,
		KeyDuration,
		KeyEncodedLibrary,
		KeyFileExtension,
		KeyFileSize,
		KeyFooterSize,
		KeyHeaderSize,
		"Comment",
		"CompleteName_Last",
		"ContentType",
		"Part",
		"FileName",
		"StreamKind",
		"CompleteName",
		"FileNameExtension",
		"Count",
		"StreamCount",
		"StreamKindID",
		"FolderName",
	},
)

var KeyGroupVideo = slices.Concat(
	KeyGroupBitRate,
	KeyGroupBitRateExtra,
	KeyGroupChroma,
	KeyGroupCommon,
	KeyGroupDimensions,
	KeyGroupFrame,
	KeyGroupFrameExtra,
	KeyGroupHDR,
	KeyGroupStream,
	KeyGroupEncoded,
	[]string{
		KeyDuration,
		KeyFormatSettingsCABAC,
		KeyFormatSettingsRefFrames,
		KeyScanType,
		"BufferSize",
		"Compression_Mode",
		"Default",
		"Delay",
		"Delay_DropFrame",
		"Delay_Settings",
		"Delay_Source",
		"DisplayAspectRatio_Original",
		"Duration_FirstFrame",
		"Forced",
		"Format_Compression",
		"Format_Settings_GOP",
		"Format_Settings_Packing",
		"Format_Settings_SliceCount",
		"Format_Tier",
		"Source_Duration",
		"Source_StreamSize",
		"Standard",
		"UniqueID",
	},
)

var KeyGroupAudio = slices.Concat(
	KeyGroupBitRate,
	KeyGroupBitRateExtra,
	KeyGroupChannelInfo,
	KeyGroupCommon,
	KeyGroupFrame,
	KeyGroupFrameExtra,
	KeyGroupSampling,
	KeyGroupSource,
	KeyGroupStream,
	KeyGroupEncoded,
	KeyGroupFormatSettings,
	[]string{
		KeyCompressionMode,
		KeyDuration,
		KeyEncodedLibrary,
		KeyFormatAdditionalFeatures,
		"AlternateGroup",
		"Forced",
		"Format_Commercial_IfAny",
		"Default",
		"Delay",
		"Delay_DropFrame",
		"Delay_Source",
		"Duration_LastFrame",
		"ServiceKind",
		"Source_Duration_LastFrame",
		"Video_Delay",
	},
)

var KeyGroupImage = slices.Concat(
	KeyGroupCommon,
	KeyGroupDimensions,
	KeyGroupChroma,
	KeyGroupHDR,
	KeyGroupStream,
	[]string{
		KeyFormatCompression,
		KeyFormatSettingsPacking,
		"Compression_Mode",
	},
)

var KeyGroupText = slices.Concat(
	KeyGroupCommon,
	[]string{},
)

var KeyGroupMenu = slices.Concat(
	KeyGroupCommon,
	[]string{
		"StreamOrder",
		"Duration",
	},
)

var KeyGroupOther = slices.Concat(
	KeyGroupCommon,
	[]string{
		"Default",
		"Duration",
		"FrameCount",
		"FrameRate",
		"FrameRate_Den",
		"FrameRate_Num",
		"StreamOrder",
		"TimeCode_FirstFrame",
		"TimeCode_LastFrame",
		"TimeCode_Stripped",
		"Type",
	},
)
