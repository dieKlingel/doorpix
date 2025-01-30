//
//  Generated code. Do not modify.
//  source: core.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use emitRequestDescriptor instead')
const EmitRequest$json = {
  '1': 'EmitRequest',
  '2': [
    {'1': 'type', '3': 1, '4': 1, '5': 9, '10': 'type'},
    {'1': 'data', '3': 2, '4': 3, '5': 11, '6': '.EmitRequest.DataEntry', '10': 'data'},
  ],
  '3': [EmitRequest_DataEntry$json],
};

@$core.Deprecated('Use emitRequestDescriptor instead')
const EmitRequest_DataEntry$json = {
  '1': 'DataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `EmitRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emitRequestDescriptor = $convert.base64Decode(
    'CgtFbWl0UmVxdWVzdBISCgR0eXBlGAEgASgJUgR0eXBlEioKBGRhdGEYAiADKAsyFi5FbWl0Um'
    'VxdWVzdC5EYXRhRW50cnlSBGRhdGEaNwoJRGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQK'
    'BXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use emitResponseDescriptor instead')
const EmitResponse$json = {
  '1': 'EmitResponse',
};

/// Descriptor for `EmitResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emitResponseDescriptor = $convert.base64Decode(
    'CgxFbWl0UmVzcG9uc2U=');

@$core.Deprecated('Use runRequestDescriptor instead')
const RunRequest$json = {
  '1': 'RunRequest',
  '2': [
    {'1': 'spec', '3': 1, '4': 3, '5': 11, '6': '.RunRequest.SpecEntry', '10': 'spec'},
    {'1': 'data', '3': 2, '4': 3, '5': 11, '6': '.RunRequest.DataEntry', '10': 'data'},
  ],
  '3': [RunRequest_SpecEntry$json, RunRequest_DataEntry$json],
};

@$core.Deprecated('Use runRequestDescriptor instead')
const RunRequest_SpecEntry$json = {
  '1': 'SpecEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use runRequestDescriptor instead')
const RunRequest_DataEntry$json = {
  '1': 'DataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `RunRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List runRequestDescriptor = $convert.base64Decode(
    'CgpSdW5SZXF1ZXN0EikKBHNwZWMYASADKAsyFS5SdW5SZXF1ZXN0LlNwZWNFbnRyeVIEc3BlYx'
    'IpCgRkYXRhGAIgAygLMhUuUnVuUmVxdWVzdC5EYXRhRW50cnlSBGRhdGEaNwoJU3BlY0VudHJ5'
    'EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaNwoJRGF0YUVudH'
    'J5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use runResponseDescriptor instead')
const RunResponse$json = {
  '1': 'RunResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'data', '3': 2, '4': 3, '5': 11, '6': '.RunResponse.DataEntry', '10': 'data'},
  ],
  '3': [RunResponse_DataEntry$json],
};

@$core.Deprecated('Use runResponseDescriptor instead')
const RunResponse_DataEntry$json = {
  '1': 'DataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `RunResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List runResponseDescriptor = $convert.base64Decode(
    'CgtSdW5SZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNzEioKBGRhdGEYAiADKAsyFi'
    '5SdW5SZXNwb25zZS5EYXRhRW50cnlSBGRhdGEaNwoJRGF0YUVudHJ5EhAKA2tleRgBIAEoCVID'
    'a2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

