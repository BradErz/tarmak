/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"fmt"
	"io"
	"net"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"

	"github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
	"github.com/jetstack/tarmak/pkg/wing/apiserver"
)

const defaultEtcdPathPrefix = "/registry/wing.tarmak.io"

type WingServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions

	StdOut io.Writer
	StdErr io.Writer
}

func NewWingServerOptions(out, errOut io.Writer) *WingServerOptions {
	o := &WingServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(defaultEtcdPathPrefix, apiserver.Scheme, apiserver.Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion)),

		StdOut: out,
		StdErr: errOut,
	}

	return o
}

// NewCommandStartMaster provides a CLI handler for 'start master' command
func NewCommandStartWingServer(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	o := NewWingServerOptions(out, errOut)

	cmd := &cobra.Command{
		Short: "Launch a wing API server",
		Long:  "Launch a wing API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunWingServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	//o.ServerRunOptions.AddUniversalFlags(flags)
	o.RecommendedOptions.Etcd.AddFlags(flags)
	o.RecommendedOptions.SecureServing.AddFlags(flags)
	//o.Admission.AddFlags(flags)

	return cmd
}

func (o *WingServerOptions) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *WingServerOptions) Complete() error {
	return nil
}

func (o WingServerOptions) Config() (*apiserver.Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewConfig(apiserver.Codecs)
	if err := o.RecommendedOptions.Etcd.ApplyTo(serverConfig); err != nil {
		return nil, err
	}
	if err := o.RecommendedOptions.SecureServing.ApplyTo(serverConfig); err != nil {
		return nil, err
	}
	/*
		if err := o.ServerRunOptions.ApplyTo(serverConfig); err != nil {
			return nil, err
		}
			if err := o.RecommendedOptions.Audit.ApplyTo(serverConfig); err != nil {
				return nil, err
			}

			if err := o.RecommendedOptions.Authentication.ApplyTo(serverConfig); err != nil {
				return nil, err
			}



			if err := o.RecommendedOptions.Features.ApplyTo(serverConfig); err != nil {
				return nil, err
			}
	*/

	config := &apiserver.Config{
		GenericConfig: serverConfig,
	}

	/*config.GenericConfig.LoopbackClientConfig = &rest.Config{
		Host: "127.0.0.1:9080",
	}

	config.GenericConfig.SecureServingInfo = &genericapiserver.SecureServingInfo{
		BindAddress: "0.0.0.0:9080",
	}

	config.GenericConfig.ApplyClientCert(
	*/

	return config, nil
}

func (o WingServerOptions) RunWingServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	/*server.GenericAPIServer.AddPostStartHook("start-sample-server-informers", func(context genericapiserver.PostStartHookContext) error {
		config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
		return nil
	})
	*/

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
