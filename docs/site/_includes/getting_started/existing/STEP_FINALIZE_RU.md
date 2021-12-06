<script type="text/javascript" src='{{ assets["getting-started.js"].digest_path }}'></script>
<script type="text/javascript" src='{{ assets["getting-started-access.js"].digest_path }}'></script>

Кластеры некоторых провайдеров могут требовать дополнительных действий, необходимых для запуска Deckhouse и его компонентов. Ниже приведены некоторые такие случаи. Если вы столкнулись с другими особенностями установки Deckhouse в существующем кластере, пожалуйста, опишите ваш опыт в [issue](https://github.com/deckhouse/deckhouse/issues). 

{% offtopic title="Кластер в VK Cloud Solutions (MailRu Cloud Solutions)" %}
<ul><li><p>В кластерах версии 1.21+ присутствует Gatekeeper (OPA), который требует выставления limits и requests подам. Pod Deckhouse архитектурно не имеет limits/requests, и, после установки Deckhouse, его Pod не будет запущен. При просмотре событий Deployment <code>deckhouse</code> вы можете увидеть ошибку:</p>
<div class="highlight"><pre><code>admission webhook "validation.gatekeeper.sh" denied the request: [container-must-have-limits] container <...> has no resource limits...</code></pre></div>

<p>Чтобы Deckhouse смог запуститься необходимо удалить Gatekeeper. Для удаления Gatekeeper выполните на узле, имеющем доступ к API кластера, следующие команды (потребуется установка <a href="https://helm.sh/" target="_blank">helm-клиента</a>):</p>
{% snippetcut selector="gatekeeper-uninstall" %}
```shell
helm delete gatekeeper --namespace opa-gatekeeper
kubectl delete crd -l gatekeeper.sh/system=yes
```
{% endsnippetcut %}
</li>
<li>Для версий кластера 1.20+ необходимо удалить также <code>metrics-server</code>. Для удаления <code>metrics-server</code> выполните на узле, имеющем доступ к API кластера, следующие команды (потребуется установка <a href="https://helm.sh/" target="_blank">helm-клиента</a>): 
{% snippetcut selector="metrics-server" %}
```shell
helm -n kube-system uninstall metrics-server
```
{% endsnippetcut %}
</li></ul>
{% endofftopic %}

<blockquote>
<p>Некоторые компоненты Deckhouse по умолчанию не работают на master-узле. Если вам необходимо разрешить компонентам Deckhouse работать на master-узле, снимите с master-узла taint, выполнив следующую команду:</p>
{% snippetcut %}
```bash
kubectl patch nodegroup master --type json -p '[{"op": "remove", "path": "/spec/nodeTemplate/taints"}]'
```
{% endsnippetcut %}
</blockquote>

Если при установке вы не включали в конфигурации Deckhouse другие модули, то единственным запущенным после установки Deckhouse модулем обладающим WEB-интерфейсом будет  модуль [внутренней документации](../../documentation/v1/modules/810-deckhouse-web/). Чтобы получить доступ к его WEB-интерфейсу или WEB-интерфейсам других модулей кластера нужно создать соответствующие DNS-записи.  

Создайте DNS-записи для доступа в WEB-интерфейсы кластера:
<ul>
  <li>Выясните публичный IP-адрес узла, на котором работает Ingress-контроллер.</li>
  <li>Если у вас есть возможность добавить DNS-запись используя DNS-сервер:
    <ul>
      <li>Если ваш шаблон DNS-имен кластера является <a href="https://en.wikipedia.org/wiki/Wildcard_DNS_record">wildcard
        DNS-шаблоном</a> (например - <code>%s.kube.my</code>), то добавьте соответствующую wildcard A-запись со значением публичного IP-адреса, который вы получили выше.
      </li>
      <li>
        Если ваш шаблон DNS-имен кластера <strong>НЕ</strong> является <a
              href="https://en.wikipedia.org/wiki/Wildcard_DNS_record">wildcard DNS-шаблоном</a> (например - <code>%s-kube.company.my</code>),
        то добавьте А или CNAME-записи со значением публичного IP-адреса, который вы
        получили выше, для следующих DNS-имен сервисов Deckhouse в вашем кластере:
        <div class="highlight">
<pre class="highlight">
<code example-hosts>dashboard.example.com
deckhouse.example.com
dex.example.com
grafana.example.com
kubeconfig.example.com
status.example.com
upmeter.example.com</code>
</pre>
        </div>
      </li>
    </ul>
  </li>

  <li><p>Если вы <strong>не</strong> имеете под управлением DNS-сервер: добавьте статические записи соответствия имен конкретных сервисов публичному IP-адресу узла, на котором работает Ingress-контроллер.</p><p>Например, на персональном Linux-компьютере, с которого необходим доступ к сервисам Deckhouse, выполните следующую команду (укажите ваш публичный IP-адрес в переменной <code>PUBLIC_IP</code>) для добавления записей в файл <code>/etc/hosts</code> (для Windows используйте файл <code>%SystemRoot%\system32\drivers\etc\hosts</code>):</p>
{% snippetcut selector="export-ip" %}
```shell
export PUBLIC_IP="<PUBLIC_IP>"
```
{% endsnippetcut %}

<p>Добавьте необходимые записи в файл <code>/etc/hosts</code>:</p>

{% snippetcut selector="example-hosts" %}
```shell
sudo -E bash -c "cat <<EOF >> /etc/hosts
$PUBLIC_IP dashboard.example.com
$PUBLIC_IP deckhouse.example.com
$PUBLIC_IP dex.example.com
$PUBLIC_IP grafana.example.com
$PUBLIC_IP kubeconfig.example.com
$PUBLIC_IP status.example.com
$PUBLIC_IP upmeter.example.com
EOF
"
```
{% endsnippetcut %}
</li></ul>
