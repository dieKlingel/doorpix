//
//  Generated code. Do not modify.
//  source: core.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'core.pb.dart' as $0;

export 'core.pb.dart';

@$pb.GrpcServiceName('Core')
class CoreClient extends $grpc.Client {
  static final _$emit = $grpc.ClientMethod<$0.EmitRequest, $0.EmitResponse>(
      '/Core/Emit',
      ($0.EmitRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.EmitResponse.fromBuffer(value));
  static final _$listen = $grpc.ClientMethod<$0.RunResponse, $0.RunRequest>(
      '/Core/Listen',
      ($0.RunResponse value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.RunRequest.fromBuffer(value));

  CoreClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options,
        interceptors: interceptors);

  $grpc.ResponseFuture<$0.EmitResponse> emit($0.EmitRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$emit, request, options: options);
  }

  $grpc.ResponseStream<$0.RunRequest> listen($async.Stream<$0.RunResponse> request, {$grpc.CallOptions? options}) {
    return $createStreamingCall(_$listen, request, options: options);
  }
}

@$pb.GrpcServiceName('Core')
abstract class CoreServiceBase extends $grpc.Service {
  $core.String get $name => 'Core';

  CoreServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.EmitRequest, $0.EmitResponse>(
        'Emit',
        emit_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.EmitRequest.fromBuffer(value),
        ($0.EmitResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RunResponse, $0.RunRequest>(
        'Listen',
        listen,
        true,
        true,
        ($core.List<$core.int> value) => $0.RunResponse.fromBuffer(value),
        ($0.RunRequest value) => value.writeToBuffer()));
  }

  $async.Future<$0.EmitResponse> emit_Pre($grpc.ServiceCall call, $async.Future<$0.EmitRequest> request) async {
    return emit(call, await request);
  }

  $async.Future<$0.EmitResponse> emit($grpc.ServiceCall call, $0.EmitRequest request);
  $async.Stream<$0.RunRequest> listen($grpc.ServiceCall call, $async.Stream<$0.RunResponse> request);
}
