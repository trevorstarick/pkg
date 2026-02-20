package mediainfo

import "slices"

var KeyGroupDates = []string{
	KeyModifiedDate,
	KeyModifiedDateLocal,
	KeyEncodedDate,
	KeyTaggedDate,
}

var KeyGroupFormat = []string{
	KeyFormat,
	KeyFormatVersion,
	KeyFormatProfile,
	KeyFormatLevel,
}

var KeyGroupFormatSettings = []string{
	KeyFormatSettingsCABAC,
	KeyFormatSettingsRefFrames,
	KeyFormatSettingsPacking,
	"Format_Settings_Endianness",
	"Format_Settings_Floor",
	"Format_Settings_Mode",
	"Format_Settings_ModeExtension",
	"Format_Settings_SBR",
	"Format_Settings_Sign",
}

var KeyGroupEncoded = []string{
	"EncodedBy",
	"Encoded_Application",
	"Encoded_Library",
	"Encoded_Library_Name",
	"Encoded_Library_Date",
	"Encoded_Library_Settings",
	"Encoded_Library_Version",
}

var KeyGroupCodec = []string{
	KeyCodecID,
	KeyCodecIDCompatible,
	"CodecID_Version",
}

var KeyGroupTrackCount = []string{
	KeyVideoCount,
	KeyAudioCount,
	KeyImageCount,
	KeyTextCount,
	"OtherCount",
}

var KeyGroupDimensions = []string{
	KeyRotation,
	KeyHeight,
	KeyWidth,
	KeySampledHeight,
	KeySampledWidth,
	KeyStoredHeight,
	KeyStoredWidth,
	KeyPixelAspectRatio,
	KeyDisplayAspectRatio,
}

var KeyGroupBitRate = slices.Concat(
	[]string{
		KeyBitRate,
		KeyBitRateMode,
		KeyBitDepth,
	},
)

var KeyGroupOverallBitRate = []string{
	KeyOverallBitRate,
	KeyOverallBitRateMode,
	"OverallBitRate_Minimum",
	"OverallBitRate_Nominal",
	"OverallBitRate_Maximum",
}

var KeyGroupBitRateExtra = slices.Concat(
	KeyGroupOverallBitRate,
	[]string{
		"BitRate_Maximum",
		"BitRate_Nominal",
		"BitRate_Minimum",
	})

var KeyGroupFrame = []string{
	KeyFrameCount,
	KeyFrameRate,
	KeyFrameRateMode,
}

var KeyGroupFrameExtra = []string{
	KeyFrameRateNum,
	KeyFrameRateDen,
	"FrameRate_Maximum",
	"FrameRate_Minimum",
	"FrameRate_Mode_Original",
	"FrameRate_Original",
}

var KeyGroupChroma = []string{
	KeyColorSpace,
	KeyChromaSubsampling,
	"ChromaSubsampling_Position",
	KeyBitDepth,
}

var KeyGroupStream = []string{
	KeyStreamSize,
	KeyStreamOrder,
	KeyIsStreamable,
}

var KeyGroupChannelInfo = []string{
	KeyChannelLayout,
	KeyChannelPositions,
	KeyChannels,
}

var KeyGroupSampling = []string{
	KeySamplesPerFrame,
	KeySamplingCount,
	KeySamplingRate,
}

var KeyGroupSource = []string{
	KeySourceDuration,
	KeySourceFrameCount,
	KeySourceStreamSize,
}

var KeyGroupHDR = []string{
	"colour_description_present",
	"colour_description_present",
	"colour_description_present_Source",
	"colour_primaries",
	"colour_primaries",
	"colour_primaries_Source",
	"colour_range",
	"colour_range",
	"colour_range_Source",
	"matrix_coefficients",
	"matrix_coefficients",
	"matrix_coefficients_Source",
	"transfer_characteristics",
	"transfer_characteristics",
	"transfer_characteristics_Source",
}

var KeyGroupTags = []string{
	"Album",
	"BPM",
	"Copyright",
	"Cover",
	"Genre",
	"Movie",
	"Performer",
	"Recorded_Date",
	"Season",
	"Track",
	"Track_Position",
}
