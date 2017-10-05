.. _design_cli_ux:

******************************
Using Tarmak command-line tool
******************************

$ tarmak

``kubectl``
###########

Run ``kubectl`` on clusters (Alias for ``$ tarmak clusters kubectl``).

Usage::

  $ tarmak kubectl

------------

``init``
########

* Initialises a provider if not existing.
* Initialises an environment if not existing.
* Initialises a cluster.

Usage::

  $ tarmak init

-------------

Resources
#########

Tarmak has 3 resources that can be acted upon - environments, providers and clusters.

Usage::

  $ tarmak [providers | environments | clusters] [command]

-------------

Providers
#########

Providers resource subcommand.

``list``
********

List providers resource.

Usage::

  $ tarmak providers list

``init``
********

Initialise providers resource.

Usage::

  $ tarmak providers init

------------

Environments
############

Environments resource subcommand.

``list``
********

List environments resource.

Usage::

  $ tarmak environments list

``init``
********

Initialise environments resource.

Usage::

  $ tarmak environments init

------------

Clusters
########

Clusters resource subcommand.

``list``
********

List clusters resource.

Usage::

  $ tarmak clusters list

``init``
********

Initialise cluster resource.

Usage::

  $ tarmak clusters init

``kubectl``
***********

Run ``kubectl`` on clusters resource.

Usage::

  $ tarmak clusters kubectl

``ssh <instance_name>``
***********************

Secure Shell into an instance on clusters.

Usage::

  $ tarmak clusters ssh <instance_name>

``apply``
*********

Apply changes to cluster (apply infrastructure changes only).

Usage::

  $ tarmak clusters apply

``plan``
********

Dry run apply.

Usage::

  $ tarmak clusters plan

``XX``
******

Does not run any infrastructure changes. Reconfigure based on configuration changes.

Usage::

  $ tarmak clusters XX

``YY``
******

Reconfigure based on infrastructure+configuration changes.

Usage::

  $ tarmak clusters YY

``YY-rolling-update``
*********************

YY with rolling update.

Usage::

  $ tarmak clusters YY-rolling-update

``instances [ list | ssh ]``
****************************

Instances on Cluster resource.

``list``
^^^^^^^^

Lists nodes of the context.

``ssh``
^^^^^^^

Alias for ``$ tarmak clusters ssh``.

Usage::

  $ tarmak clusters instances [list | ssh]

``server-pools [ list ]``
*************************

``list``
^^^^^^^^

List server pools on Cluster resource.

Usage::

  $ tarmak clusters server-pools list

``images [ list | build ]``
***************************

``list``
^^^^^^^^

List images on Cluster resource.

``build``
^^^^^^^^^

Build images of Cluster resource.

Usage::

  $ tarmak clusters images [list | build]

``debug [ terraform shell | puppet | etcd | vault ]``
*****************************************************

Used for debugging.

``terraform shell``
^^^^^^^^^^^^^^^^^^^

Debug terraform via shell.

Usage::

  $ tarmak clusters debug terraform shell

``puppet``
^^^^^^^^^^

Debug puppet.

Usage::

  $ tarmak clusters debug puppet

``etcd``
^^^^^^^^

Debug etcd.

Usage::

  $ tarmak clusters debug etcd

``vault``
^^^^^^^^^

Debug vault.

Usage::

  $ tarmak clusters debug vault

------------

Relationships
#############

The relationship between Providers, Environments and Clusters is as follows:

Provider (many) -> Environment (one)

Environment (many) -> Cluster (one)

Changed Names
#############

+-----------+-------------+
| Old Name  | New Name    |
+===========+=============+
| NodeGroup | Server Pool |
+-----------+-------------+
| Context   | Cluster     |
+-----------+-------------+
|  Nodes    | Instances   |
+-----------+-------------+
