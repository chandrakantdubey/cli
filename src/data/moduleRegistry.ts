import intro from './modules/1-intro.json';
import mentalModel from './modules/2-mental-model.json';
import coreTerms from './modules/3-core-terms.json';
import filesystem from './modules/4-filesystem.json';
import users from './modules/5-users.json';
import permissions from './modules/6-permissions.json';
import fileOps from './modules/7-file-ops.json';
import ioRedirection from './modules/8-io-redirection.json';
import textManipulation from './modules/9-text-manipulation.json';
import findingThings from './modules/10-finding-things.json';
import envVars from './modules/11-env-vars.json';
import packages from './modules/12-packages.json';
import networking from './modules/13-networking.json';
import ssh from './modules/14-ssh.json';
import processMgmt from './modules/15-process-mgmt.json';
import shellScripting from './modules/16-shell-scripting.json';
import git from './modules/17-git.json';
import archiving from './modules/18-archiving.json';
import textAdvanced from './modules/19-text-advanced.json';
import nodeCli from './modules/20-node-cli.json';
import pythonCli from './modules/21-python-cli.json';
import customization from './modules/22-customization.json';
import shellExpansion from './modules/23-shell-expansion.json';
import storageFilesystems from './modules/24-storage-filesystems.json';
import systemdServices from './modules/25-systemd-services.json';
import linuxTroubleshooting from './modules/26-linux-troubleshooting.json';
import linuxSecurity from './modules/27-linux-security.json';
import containersCli from './modules/28-containers-cli.json';
import observabilityCli from './modules/29-observability-cli.json';
import productionOperations from './modules/30-production-operations.json';
import automationScripting from './modules/31-automation-and-scripting.json';
import linuxCloudNative from './modules/32-linux-for-cloud-native-engineers.json';
import shellPortability from './modules/33-shell-portability.json';
import capstoneLinuxOperations from './modules/34-capstone-linux-operations.json';

export interface Topic {
    id: string;
    title: string;
    level: 'Basic' | 'Advanced';
    content: string;
}

export const courseContent: Topic[] = [
    intro,
    mentalModel,
    coreTerms,
    filesystem,
    users,
    permissions,
    fileOps,
    ioRedirection,
    textManipulation,
    findingThings,
    envVars,
    packages,
    networking,
    ssh,
    processMgmt,
    shellScripting,
    git,
    archiving,
    textAdvanced,
    nodeCli,
    pythonCli,
    customization,
    shellExpansion,
    storageFilesystems,
    systemdServices,
    linuxTroubleshooting,
    linuxSecurity,
    containersCli,
    observabilityCli,
    productionOperations,
    automationScripting,
    linuxCloudNative,
    shellPortability,
    capstoneLinuxOperations
] as Topic[];
