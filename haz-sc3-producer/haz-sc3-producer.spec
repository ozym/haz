# Disable the stupid stuff rpm distros include in the build process by default:
#   Disable any prep shell actions. replace them with simply 'true'
%define __spec_prep_post true
%define __spec_prep_pre true
#   Disable any build shell actions. replace them with simply 'true'
%define __spec_build_post true
%define __spec_build_pre true
#   Disable any install shell actions. replace them with simply 'true'
%define __spec_install_post true
%define __spec_install_pre true
#   Disable any clean shell actions. replace them with simply 'true'
%define __spec_clean_post true
%define __spec_clean_pre true
# Disable checking for unpackaged files ?
#%undefine __check_files

%define debug_package   %{nil}

%if 0%{!?rev:1}
%define rev             %(git rev-parse HEAD)
%endif
%define shortrev        %(r=%{rev}; echo ${r:0:7})

%define gh_user         GeoNet
%define gh_name         haz-sc3-producer
%define gh_tar          %{gh_user}-%{gh_name}-%{shortrev}
%define import_path     github.com/%{gh_user}/%{gh_name}

Name:       haz-sc3-producer
Version:    0.1
Release:    %{?rel}git%{shortrev}%{?dist}
Summary:    Sends SeisComPML file to AWS SNS as Haz JSON messages.

Group:		Applications/Webapps
License:	GNS
URL:		https://%{import_path}
Source0:	https://%{import_path}/tarball/master/%{gh_tar}.tar.gz

BuildRequires:	golang

%description
Sends SeisComPML file to AWS SNS as Haz JSON messages.

%prep
# noop

%build
# noop

%install
# noop

%clean
# noop


%pre
getent group alert >/dev/null || groupadd -g 3506 alert
getent passwd alert >/dev/null || useradd -d /home/alert -u 3506 -g alert alert


%post
/sbin/chkconfig --add haz-sc3-producer


%preun
# Checks that this is the actual deinstallation of the package, as opposed
# to just removing the old package on upgrade.
if [ $1 = 0 ] ; then
    /sbin/service haz-sc3-producer stop >/dev/null 2>&1
    /sbin/chkconfig --del haz-sc3-producer
fi


%postun
# Checks that this is an upgrade of the package.
if [ $1 -ge 1 ] ; then
    /sbin/service haz-sc3-producer condrestart >/dev/null 2>&1 || :
fi


%files
%defattr(-,root,root,-)
%doc README.md
%config(noreplace) %{_sysconfdir}/sysconfig/haz-sc3-producer.json
%attr(755,root,root) %{_bindir}/haz-sc3-producer
%attr(755,root,root) %{_initrddir}/haz-sc3-producer


%changelog
