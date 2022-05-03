import 'package:dmch_gui/bloc/auth.dart';
import 'package:dmch_gui/views/login.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:provider/provider.dart';

import '../api/dmapi.dart';
import 'media_view.dart';

// class MainView extends StatefulWidget {
//   const MainView({Key? key}) : super(key: key);

//   @override
//   State<MainView> createState() => _MainViewState();
// }

class MainView extends StatelessWidget {
  const MainView({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Provider<DmApiClient>(
      create: (context) => DmApiClient(),
      child: BlocProvider(
        create: (context) => AuthenticationBloc(
          apiClient: Provider.of<DmApiClient>(context, listen: false),
        ),
        child: BlocBuilder<AuthenticationBloc, AuthenticationState>(
          builder: (context, state) {
            switch (state.status) {
              case AuthenticationStatus.unauthenticated:
                return const LoginView();
              case AuthenticationStatus.authenticated:
                return const MediaView();
              default:
                return const Center(child: CircularProgressIndicator());
            }
          },
        ),
      ),
    );
  }
}
