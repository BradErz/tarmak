require 'spec_helper'

describe 'prometheus::prometheus_scrape_config', :type => :define do
  let(:pre_condition) {[
    'include kubernetes::apiserver',
  ]}

  context 'test scrape static_configs definition' do
    let(:title) do
      'etcd_k8s'
    end

    let(:etcd_cluster) { ['192.168.1.2', '192.168.1.3'] }
    let :params do
      {
        :static_configs    => { "- targets" => "inline_template(<%= @etcd_cluster.join(',') %>)"},
        :order             => '02',
      }
    end

    it do
      should contain_concat__fragment("kubectl-apply-prometheus-scrape-config-etcd_k8s")
        .with_content(/- targets: [ '192.168.1.2' ]/)
    end
  end

  context 'test scrape kubernetes_sd_configs definition' do
    let(:title) do
      'kubernetes-apiservers'
    end
    let :params do
      {
        :kubernetes_sd_configs => { "- role" => "endpoints"},
        :order                 => '02',
      }
    end
    it do
      should contain_concat__fragment("kubectl-apply-prometheus-scrape-config-kubernetes-apiservers")
        .with_content(/- role: endpoints/)
    end
  end
end
