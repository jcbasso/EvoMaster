package com.foo.rpc.grpc.examples.spring.branches.generated;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.57.2)",
    comments = "Source: branches.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class BranchGRPCServiceGrpc {

  private BranchGRPCServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "BranchGRPCService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
      com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getPosMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "pos",
      requestType = com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto.class,
      responseType = com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
      com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getPosMethod() {
    io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto, com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getPosMethod;
    if ((getPosMethod = BranchGRPCServiceGrpc.getPosMethod) == null) {
      synchronized (BranchGRPCServiceGrpc.class) {
        if ((getPosMethod = BranchGRPCServiceGrpc.getPosMethod) == null) {
          BranchGRPCServiceGrpc.getPosMethod = getPosMethod =
              io.grpc.MethodDescriptor.<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto, com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "pos"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto.getDefaultInstance()))
              .setSchemaDescriptor(new BranchGRPCServiceMethodDescriptorSupplier("pos"))
              .build();
        }
      }
    }
    return getPosMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
      com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getNegMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "neg",
      requestType = com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto.class,
      responseType = com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
      com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getNegMethod() {
    io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto, com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getNegMethod;
    if ((getNegMethod = BranchGRPCServiceGrpc.getNegMethod) == null) {
      synchronized (BranchGRPCServiceGrpc.class) {
        if ((getNegMethod = BranchGRPCServiceGrpc.getNegMethod) == null) {
          BranchGRPCServiceGrpc.getNegMethod = getNegMethod =
              io.grpc.MethodDescriptor.<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto, com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "neg"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto.getDefaultInstance()))
              .setSchemaDescriptor(new BranchGRPCServiceMethodDescriptorSupplier("neg"))
              .build();
        }
      }
    }
    return getNegMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
      com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getEqMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "eq",
      requestType = com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto.class,
      responseType = com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
      com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getEqMethod() {
    io.grpc.MethodDescriptor<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto, com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> getEqMethod;
    if ((getEqMethod = BranchGRPCServiceGrpc.getEqMethod) == null) {
      synchronized (BranchGRPCServiceGrpc.class) {
        if ((getEqMethod = BranchGRPCServiceGrpc.getEqMethod) == null) {
          BranchGRPCServiceGrpc.getEqMethod = getEqMethod =
              io.grpc.MethodDescriptor.<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto, com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "eq"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto.getDefaultInstance()))
              .setSchemaDescriptor(new BranchGRPCServiceMethodDescriptorSupplier("eq"))
              .build();
        }
      }
    }
    return getEqMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static BranchGRPCServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BranchGRPCServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BranchGRPCServiceStub>() {
        @java.lang.Override
        public BranchGRPCServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BranchGRPCServiceStub(channel, callOptions);
        }
      };
    return BranchGRPCServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BranchGRPCServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BranchGRPCServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BranchGRPCServiceBlockingStub>() {
        @java.lang.Override
        public BranchGRPCServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BranchGRPCServiceBlockingStub(channel, callOptions);
        }
      };
    return BranchGRPCServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static BranchGRPCServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BranchGRPCServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BranchGRPCServiceFutureStub>() {
        @java.lang.Override
        public BranchGRPCServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BranchGRPCServiceFutureStub(channel, callOptions);
        }
      };
    return BranchGRPCServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public interface AsyncService {

    /**
     */
    default void pos(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request,
        io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPosMethod(), responseObserver);
    }

    /**
     */
    default void neg(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request,
        io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getNegMethod(), responseObserver);
    }

    /**
     */
    default void eq(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request,
        io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getEqMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service BranchGRPCService.
   */
  public static abstract class BranchGRPCServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return BranchGRPCServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service BranchGRPCService.
   */
  public static final class BranchGRPCServiceStub
      extends io.grpc.stub.AbstractAsyncStub<BranchGRPCServiceStub> {
    private BranchGRPCServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BranchGRPCServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BranchGRPCServiceStub(channel, callOptions);
    }

    /**
     */
    public void pos(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request,
        io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPosMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void neg(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request,
        io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getNegMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void eq(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request,
        io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getEqMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service BranchGRPCService.
   */
  public static final class BranchGRPCServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<BranchGRPCServiceBlockingStub> {
    private BranchGRPCServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BranchGRPCServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BranchGRPCServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto pos(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPosMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto neg(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getNegMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto eq(com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getEqMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service BranchGRPCService.
   */
  public static final class BranchGRPCServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<BranchGRPCServiceFutureStub> {
    private BranchGRPCServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BranchGRPCServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BranchGRPCServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> pos(
        com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPosMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> neg(
        com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getNegMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto> eq(
        com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getEqMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_POS = 0;
  private static final int METHODID_NEG = 1;
  private static final int METHODID_EQ = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AsyncService serviceImpl;
    private final int methodId;

    MethodHandlers(AsyncService serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_POS:
          serviceImpl.pos((com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto) request,
              (io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>) responseObserver);
          break;
        case METHODID_NEG:
          serviceImpl.neg((com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto) request,
              (io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>) responseObserver);
          break;
        case METHODID_EQ:
          serviceImpl.eq((com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto) request,
              (io.grpc.stub.StreamObserver<com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  public static final io.grpc.ServerServiceDefinition bindService(AsyncService service) {
    return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
        .addMethod(
          getPosMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
              com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>(
                service, METHODID_POS)))
        .addMethod(
          getNegMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
              com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>(
                service, METHODID_NEG)))
        .addMethod(
          getEqMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.foo.rpc.grpc.examples.spring.branches.generated.BranchesPostDto,
              com.foo.rpc.grpc.examples.spring.branches.generated.BranchesResponseDto>(
                service, METHODID_EQ)))
        .build();
  }

  private static abstract class BranchGRPCServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BranchGRPCServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.foo.rpc.grpc.examples.spring.branches.generated.Branches.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("BranchGRPCService");
    }
  }

  private static final class BranchGRPCServiceFileDescriptorSupplier
      extends BranchGRPCServiceBaseDescriptorSupplier {
    BranchGRPCServiceFileDescriptorSupplier() {}
  }

  private static final class BranchGRPCServiceMethodDescriptorSupplier
      extends BranchGRPCServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    BranchGRPCServiceMethodDescriptorSupplier(java.lang.String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (BranchGRPCServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new BranchGRPCServiceFileDescriptorSupplier())
              .addMethod(getPosMethod())
              .addMethod(getNegMethod())
              .addMethod(getEqMethod())
              .build();
        }
      }
    }
    return result;
  }
}
