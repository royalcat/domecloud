import 'package:bloc/bloc.dart';
import 'package:dmch_gui/api/dmapi.dart';
import 'package:dmch_gui/models/users.dart';
import 'package:equatable/equatable.dart';

class AuthenticationBloc extends Bloc<AuthenticationEvent, AuthenticationState> {
  final DmApiClient apiClient;

  AuthenticationBloc({
    required this.apiClient,
  }) : super(const AuthenticationState.unauthenticated()) {
    on<AuthenticationLogInEvent>(_onAuthenticationLogIn);
    on<AuthenticationLogoutEvent>(_onAuthenticationLogOut);
  }

  void _onAuthenticationLogIn(
    AuthenticationLogInEvent event,
    Emitter<AuthenticationState> emit,
  ) async {
    switch (state.status) {
      case AuthenticationStatus.authenticated:
      case AuthenticationStatus.unauthenticated:
        emit(const AuthenticationState.unknown());
        final user = await apiClient.logIn(
          username: event.username,
          password: event.password,
        );
        return emit(user != null
            ? AuthenticationState.authenticated(user)
            : const AuthenticationState.unauthenticated());
      case AuthenticationStatus.unknown:
        break;
    }
  }

  void _onAuthenticationLogOut(
    AuthenticationLogoutEvent event,
    Emitter<AuthenticationState> emit,
  ) {
    apiClient.logOut();
  }
}

abstract class AuthenticationEvent extends Equatable {
  const AuthenticationEvent();

  @override
  List<Object> get props => [];
}

class AuthenticationLogInEvent extends AuthenticationEvent {
  final String username;
  final String password;

  const AuthenticationLogInEvent(this.username, this.password);

  @override
  List<Object> get props => [username, password];
}

class AuthenticationLogoutEvent extends AuthenticationEvent {}

enum AuthenticationStatus { unknown, authenticated, unauthenticated }

class AuthenticationState extends Equatable {
  const AuthenticationState._({
    this.status = AuthenticationStatus.unknown,
    this.user = const User(username: "", password: ""),
  });

  const AuthenticationState.unknown() : this._();

  const AuthenticationState.authenticated(User user)
      : this._(status: AuthenticationStatus.authenticated, user: user);

  const AuthenticationState.unauthenticated()
      : this._(status: AuthenticationStatus.unauthenticated);

  final AuthenticationStatus status;
  final User user;

  @override
  List<Object> get props => [status, user];
}
