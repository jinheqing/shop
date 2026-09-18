import 'package:dio/dio.dart';

class ApiService {
  static final ApiService _i = ApiService._();
  factory ApiService() => _i;
  ApiService._() {
    _dio = Dio(BaseOptions(baseUrl: 'http://10.0.2.2:8080/api/v1', connectTimeout: Duration(seconds: 10), receiveTimeout: Duration(seconds: 15)));
  }
  late final Dio _dio;
  String? _token;

  void setToken(String t) => _token = t;

  Future<Map<String, dynamic>> login(String email, String password) async {
    final r = await _dio.post('/staff/login', data: {'email': email, 'password': password});
    _token = r.data['access_token'];
    return r.data;
  }

  Future<List<dynamic>> getOrders() async => (await _dio.get('/orders', options: _auth())).data['items'];
  Future<Map<String, dynamic>> getOrder(int id) async => (await _dio.get('/orders/$id', options: _auth())).data;
  Future<Map<String, dynamic>> updateOrderState(int id, String state) async =>
      (await _dio.post('/orders/$id/state', data: {'state': state}, options: _auth())).data;

  Future<List<dynamic>> getLiveRooms() async => (await _dio.get('/live-rooms', options: _auth())).data['items'];
  Future<List<dynamic>> getNodes() async => (await _dio.get('/nodes', options: _auth())).data['items'];
  Future<List<dynamic>> getCustomProducts() async => (await _dio.get('/custom-products', options: _auth())).data['items'];
  Future<List<dynamic>> getSlowPresets() async => (await _dio.get('/slow-presets', options: _auth())).data;

  Options _auth() => Options(headers: {'Authorization': 'Bearer $_token'});
}
