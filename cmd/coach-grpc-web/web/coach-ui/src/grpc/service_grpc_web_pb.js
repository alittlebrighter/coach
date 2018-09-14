/**
 * @fileoverview gRPC-Web generated client stub for coach
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');


var coach_pb = require('./coach_pb.js')
const proto = {};
proto.coach = require('./service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.coach.CoachRPCClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.coach.CoachRPCPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!proto.coach.CoachRPCClient} The delegate callback based client
   */
  this.delegateClient_ = new proto.coach.CoachRPCClient(
      hostname, credentials, options);

};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.coach.ScriptsQuery,
 *   !proto.coach.GetScriptsResponse>}
 */
const methodInfo_Scripts = new grpc.web.AbstractClientBase.MethodInfo(
  proto.coach.GetScriptsResponse,
  /** @param {!proto.coach.ScriptsQuery} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.coach.GetScriptsResponse.deserializeBinary
);


/**
 * @param {!proto.coach.ScriptsQuery} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.coach.GetScriptsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.coach.GetScriptsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.coach.CoachRPCClient.prototype.scripts =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/coach.CoachRPC/Scripts',
      request,
      metadata,
      methodInfo_Scripts,
      callback);
};


/**
 * @param {!proto.coach.ScriptsQuery} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.coach.GetScriptsResponse>}
 *     The XHR Node Readable Stream
 */
proto.coach.CoachRPCPromiseClient.prototype.scripts =
    function(request, metadata) {
  return new Promise((resolve, reject) => {
    this.delegateClient_.scripts(
      request, metadata, (error, response) => {
        error ? reject(error) : resolve(response);
      });
  });
};


module.exports = proto.coach;

